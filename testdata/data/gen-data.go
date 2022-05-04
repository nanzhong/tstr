package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/nanzhong/tstr/db"

	"github.com/sethvargo/go-diceware/diceware"
)

const (
	numRunners     = 25
	numTests       = 25
	numLabels      = 10
	numRunsPerTest = 30
)

type Runner struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	AcceptTestLabels map[string]string  `json:"accept_test_labels"`
	RejectTestLabels map[string]string  `json:"reject_test_labels"`
	RegisteredAt     pgtype.Timestamptz `json:"registered_at"`
	ApprovedAt       pgtype.Timestamptz `json:"approved_at"`
	RevokedAt        pgtype.Timestamptz `json:"revoked_at"`
	LastHeartbeatAt  pgtype.Timestamptz `json:"last_heartbeat_at"`
}

type Test struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Labels       map[string]string  `json:"labels"`
	CronSchedule string             `json:"cron_schedule"`
	RegisteredAt pgtype.Timestamptz `json:"registered_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
	ArchivedAt   pgtype.Timestamptz `json:"archived_at"`
}

type TestRunConfig struct {
	ID             string             `json:"id"`
	TestID         string             `json:"test_id"`
	ContainerImage string             `json:"container_image"`
	Command        string             `json:"command"`
	Args           []string           `json:"args"`
	Env            map[string]string  `json:"env"`
	CreatedAt      pgtype.Timestamptz `json:"created_at"`
}

type RunLogEntry struct {
	Time       string `json:"time"`
	OutputType int    `json:"output_type"`
	Data       []byte `json:"data"`
}

type Run struct {
	ID              string             `json:"id"`
	TestID          string             `json:"test_id"`
	TestRunConfigID string             `json:"test_run_config_id"`
	RunnerID        string             `json:"runner_id"`
	Result          db.RunResult       `json:"result"`
	Logs            []RunLogEntry      `json:"logs"`
	ScheduledAt     pgtype.Timestamptz `json:"scheduled_at"`
	StartedAt       pgtype.Timestamptz `json:"started_at"`
	FinishedAt      pgtype.Timestamptz `json:"finished_at"`
}

func genRandomName(numWords int) string {
	words, err := diceware.Generate(numWords)
	panicIf(err)
	return strings.Join(words, " ")
}

func main() {
	rand.Seed(21)

	presetLabels := map[string][]string{}

	// for each label tag generates 3 values
	for i := 0; i < numLabels; i++ {
		label := genRandomName(1)
		labelValues, err := diceware.Generate(3)
		panicIf(err)
		presetLabels[label] = labelValues
	}

	randomLabels := func(numLabels int) map[string]string {
		labels := map[string]string{}

		i := 0
		for k, v := range presetLabels {
			labels[k] = v[rand.Intn(len(v))]
			i++
			if i == numLabels {
				break
			}
		}

		return labels
	}

	randomCronSchedule := func() string {
		cronStrings := []string{
			"0 8 * * *",
			"30 12 * * 3",
			"* * * * *",
			"5 0 * 8 *",
			"15 14 1 * *",
			"0 22 * * 1-5",
			"23 0-20/2 * * *",
			"5 4 * * sun",
		}
		return cronStrings[rand.Intn(len(cronStrings))]
	}

	runners := []Runner{}

	for i := 0; i < numRunners; i++ {
		// going simple for now:
		//  - approval a few minutes after registration
		//  - last heartbeat a few minutes after approval
		//
		registeredAt := time.Now().Add(time.Duration(-4*24) * time.Hour).Add(time.Hour * time.Duration(rand.Intn(12)))
		approvedAt := registeredAt.Add(time.Minute * time.Duration(rand.Intn(60)))
		lastHeartbeatAt := approvedAt.Add(time.Minute * time.Duration(rand.Intn(30)))

		runner := Runner{
			ID:               uuid.New().String(),
			Name:             fmt.Sprintf("runner-%d", i),
			AcceptTestLabels: map[string]string{},
			RejectTestLabels: map[string]string{},
			RegisteredAt:     pgtype.Timestamptz{Time: registeredAt, Status: pgtype.Present},
			ApprovedAt:       pgtype.Timestamptz{Time: approvedAt, Status: pgtype.Present},
			RevokedAt:        pgtype.Timestamptz{Status: pgtype.Null},
			LastHeartbeatAt:  pgtype.Timestamptz{Status: pgtype.Present, Time: lastHeartbeatAt},
		}
		runners = append(runners, runner)
	}

	randomRunner := func() Runner {
		return runners[rand.Intn(len(runners))]
	}

	tests := []Test{}
	testRunConfigs := []TestRunConfig{}
	testRuns := []Run{}

	for i := 0; i < numTests; i++ {

		// registered in the last [5:35[ days
		registeredAt := time.Now().Add(-time.Hour * 24 * time.Duration(5+rand.Intn(30)))

		archivedAt := pgtype.Timestamptz{}
		if rand.Float64() <= 0.15 {
			archivedAt.Status = pgtype.Present
			archivedAt.Time = time.Now().Add(-time.Duration(rand.Intn(5)) * time.Hour)
		} else {
			archivedAt.Status = pgtype.Null
		}

		test := Test{
			ID:           uuid.New().String(),
			Name:         genRandomName(3),
			Labels:       randomLabels(2),
			CronSchedule: randomCronSchedule(),
			RegisteredAt: pgtype.Timestamptz{Time: registeredAt, Status: pgtype.Present},
			UpdatedAt:    pgtype.Timestamptz{Time: registeredAt, Status: pgtype.Present},
			ArchivedAt:   archivedAt,
		}

		// for each test create 1 (for now) TestRunConfig

		testRunConfig := TestRunConfig{
			ID:             uuid.New().String(),
			TestID:         test.ID,
			ContainerImage: "busybox:latest",
			Command:        "/bin/sh",
			Args:           []string{"-c", "echo hello test $TSTR_TEST_ID"},
			Env:            map[string]string{"LANG": "C"},
			CreatedAt:      pgtype.Timestamptz{Time: test.UpdatedAt.Time.Add(1 + time.Duration(rand.Intn(10))*time.Minute), Status: pgtype.Present},
		}
		testRunConfigs = append(testRunConfigs, testRunConfig)
		tests = append(tests, test)

		medianPassDurationSeconds := 300
		medianFailDurationSeconds := 100

		// 1 test => 1 testrunconfig => $numRunsPerTest runs
		for r := 0; r < numRunsPerTest; r++ {

			result := db.RunResultPass
			medianDuration := medianPassDurationSeconds

			if rand.Float64() < 0.10 {
				result = db.RunResultFail
				medianDuration = medianFailDurationSeconds
			}

			scheduledAt := testRunConfig.CreatedAt.Time.Add(time.Duration(rand.Intn(10)) * time.Minute)
			startedAt := scheduledAt.Add(time.Duration(rand.Intn(30)) * time.Second)
			finishedAt := startedAt.Add(time.Second * time.Duration(float64(medianDuration)+rand.NormFloat64()*40))

			run := Run{
				ID:              uuid.New().String(),
				TestID:          test.ID,
				TestRunConfigID: testRunConfig.ID,
				RunnerID:        randomRunner().ID,
				Result:          result,
				Logs:            []RunLogEntry{},
				ScheduledAt:     pgtype.Timestamptz{Status: pgtype.Present, Time: scheduledAt},
				StartedAt:       pgtype.Timestamptz{Status: pgtype.Present, Time: startedAt},
				FinishedAt:      pgtype.Timestamptz{Status: pgtype.Present, Time: finishedAt},
			}

			testRuns = append(testRuns, run)
		}
	}

	dumpIntoJsonl(anyToSlice(runners), "runners.json")
	dumpIntoJsonl(anyToSlice(tests), "tests.json")
	dumpIntoJsonl(anyToSlice(testRunConfigs), "test_run_configs.json")
	dumpIntoJsonl(anyToSlice(testRuns), "runs.json")
}

func dumpIntoJsonl(data any, outputFile string) {
	f, err := os.Create(outputFile)
	panicIf(err)
	defer f.Close()

	switch d := data.(type) {
	case []any:
		for _, row := range d {
			jsonLine, err := json.Marshal(row)
			panicIf(err)
			f.Write(jsonLine)
			f.Write([]byte("\n"))
		}
	default:
		log.Panicf("dunno what to do with %T\n", data)
	}
}

func anyToSlice(slice any) []any {
	sliceVal := reflect.ValueOf(slice)
	if sliceVal.Kind() != reflect.Slice {
		log.Panicln("dunno what to do with %T\n", slice)
	}

	if sliceVal.IsNil() {
		return nil
	}

	ret := make([]any, sliceVal.Len())

	for i := 0; i < sliceVal.Len(); i++ {
		ret[i] = sliceVal.Index(i).Interface()
	}

	return ret
}

func panicIf(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
