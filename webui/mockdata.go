package webui

// func genRuns() []common.Run {
// 	n := 65
// 	runs := make([]common.Run, n+1)
//
// 	latestTimestamp := time.Now()
//
// 	for i := n; i >= 0; i-- {
// 		scheduledAt := latestTimestamp.Add(-30 * time.Minute).Add(time.Duration(rand.Intn(100)) * time.Second)
// 		startedAt := scheduledAt.Add(time.Second * (10 + time.Duration(rand.Intn(30))))
// 		latestTimestamp = scheduledAt
//
// 		runs[i].ScheduledAt = timestamppb.New(scheduledAt)
// 		// a run started between 10 and 40 secs after being scheduled
// 		runs[i].StartedAt = timestamppb.New(startedAt)
//
// 		// a run took between 5 and 8 minutes
// 		runs[i].FinishedAt = timestamppb.New(startedAt.Add(time.Minute * (5 + time.Duration(rand.Intn(4)))))
//
// 		runs[i].Id = uuid.New().String()
// 		runs[i].TestId = uuid.New().String()
// 		runs[i].RunnerId = fmt.Sprintf("runner-%d", rand.Intn(16))
//
// 		runs[i].TestRunConfig = &common.Test_RunConfig{Args: []string{"-e=mc2"}}
// 	}
//
// 	// the last 2 runs are pending
// 	runs[n].FinishedAt = nil
// 	runs[n-1].FinishedAt = nil
//
// 	// the last one is yet to start
// 	runs[n].StartedAt = nil
//
// 	return runs
// }
