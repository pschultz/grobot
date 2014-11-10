package grobot

import "github.com/fgrosse/grobot/log"

type HookType int

const (
	HookBefore = iota
	HookAfter  = iota
)

type TaskHook struct {
	Typ     HookType
	SubTask string
}

func (h *TaskHook) String() string {
	switch h.Typ {
	case HookBefore:
		return "BEFORE"
	default:
		return "UNKOWN HOOK"
	}
}

var hooks = map[string][]*TaskHook{}

func RegisterTaskHook(hookType HookType, parentTaskName string, subTask string) {
	hook := TaskHook{hookType, subTask}
	hooks[parentTaskName] = append(hooks[parentTaskName], &hook)
}

// TODO return whether the hooked task had an update
// TODO check for recursions
func checkHooks(hookType HookType, invokedName string, recursionDepth int) (bool, error) {
	var taskHooks []*TaskHook
	var hookRegistered bool
	taskHooks, hookRegistered = hooks[invokedName]
	if hookRegistered == false {
		return false, nil
	}

	returnValues := make([]bool, len(taskHooks))
	for i, hook := range taskHooks {
		if hook.Typ != hookType {
			return false, nil
		}

		log.Debug("Invoking hook %s for target [<strong>%s</strong>]", hook.Typ, invokedName)
		wasUpdated, err := InvokeTask(hook.SubTask, recursionDepth+1)
		if err != nil {
			return false, err
		}
		returnValues[i] = wasUpdated
	}

	for _, wasUpdated := range returnValues {
		if wasUpdated {
			return true, nil
		}
	}
	return false, nil
}
