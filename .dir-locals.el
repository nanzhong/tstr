;;; Directory Local Variables            -*- no-byte-compile: t -*-
;;; For more information see (info "(emacs) Directory Variables")

((go-mode . ((eglot-workspace-configuration . (:gopls (:buildFlags ["-tags=integration,ui"])))))
 (magit-status-mode . ((magit-todos-exclude-globs . (".git/" "vendor/")))))
