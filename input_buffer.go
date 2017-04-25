package main

type InputBuffer struct {
	buffer      []rune
	keyBindings KeyBindings
}

func NewInputBuffer(keyBindings KeyBindings) *InputBuffer {
	return &InputBuffer{
		buffer:      make([]rune, 0),
		keyBindings: keyBindings,
	}
}

func (inputBuffer *InputBuffer) Append(input string) {
	inputBuffer.buffer = append(inputBuffer.buffer, []rune(input)...)
}

func (inputBuffer *InputBuffer) prepend(keystring string) {
	inputBuffer.buffer = append([]rune(keystring), inputBuffer.buffer...)
}

func (inputBuffer *InputBuffer) pop() (char rune) {
	char = inputBuffer.buffer[0]
	inputBuffer.buffer = inputBuffer.buffer[1:]
	return
}

func (inputBuffer *InputBuffer) hasInput() bool {
	return len(inputBuffer.buffer) > 0
}

func (inputBuffer *InputBuffer) Process(viewHierarchy ViewHierarchy) (action Action, keystring string) {
	if !inputBuffer.hasInput() {
		return
	}

	keyBuffer := make([]rune, 0)
	keyBindings := inputBuffer.keyBindings
	isPrefix := false

OuterLoop:
	for inputBuffer.hasInput() {
		keyBuffer = append(keyBuffer, inputBuffer.pop())
		binding, prefix := keyBindings.Binding(viewHierarchy, string(keyBuffer))

		switch {
		case prefix:
			if len(inputBuffer.buffer) == 0 {
				inputBuffer.prepend(string(keyBuffer))
				return
			} else {
				isPrefix = true
			}
		case binding.bindingType == BT_ACTION:
			if binding.action != ACTION_NONE {
				action = binding.action
			} else if isPrefix {
				inputBuffer.prepend(string(keyBuffer[1:]))
				keyBuffer = keyBuffer[0:1]
			}

			break OuterLoop
		case binding.bindingType == BT_KEYSTRING:
			inputBuffer.prepend(binding.keystring)
			keyBuffer = keyBuffer[0:0]
			isPrefix = false
		}
	}

	keystring = string(keyBuffer)

	return
}