(defgroup hound nil
  "Variables related to hound."
  :prefix "hound-"
  :group 'tools)

(defcustom hound-command "hound"
  "The name of the hound program."
  :type '(string)
  :group 'hound)

;;;###autoload
(defun hound-search (pattern)
  "Search files in the index for PATTERN."
  (interactive
   (list
    (read-string "Pattern: " (thing-at-point 'symbol))))
  (compilation-start
   (concat hound-command " -like-grep " pattern)
   'grep-mode)
  )

(provide 'hound)
