package painter

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"golang.org/x/exp/shiny/screen"
)

func TestLoop_Start(t *testing.T) {
	screen := &mockScreen{}
	loop := &Loop{
		Receiver: &testReceiver{},
	}
	loop.Start(screen)

	if loop.next == nil || loop.prev == nil {
		t.Error("unexpected nil texture")
	}

	loop.StopAndWait()
}

func TestLoop_Post(t *testing.T) {
	var (
		loop Loop
		tr   testReceiver
	)

	loop.Receiver = &tr

	loop.Start(mockScreen{})
	loop.Post(OperationFunc(WhiteFill))
	loop.Post(OperationFunc(GreenFill))
	loop.Post(UpdateOp)
	if tr.LastTexture != nil {
		t.Fatal("Reciever got the texture too early")
	}
	loop.StopAndWait()

	tx, ok := tr.LastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Reciever still has not texture")
	}
	if tx.FillCnt != 2 {
		t.Error("Unexpected number of fill calls:", tx.FillCnt)
	}

}

func TestMessageQueue_push_pull_empty(t *testing.T) {
	msgQue := &messageQueue{}
	operation := &testOperation{}

	msgQue.push(operation)

	if len(msgQue.ops) != 1 || msgQue.ops[0] != operation {
		t.Error("failed to push operation into queue")
	}

	pulledOp := msgQue.pull()
	if pulledOp != operation {
		t.Error("failed to pull operation from queue")
	}

	if !msgQue.empty() {
		t.Error("expected queue to be empty")
	}
}

func TestMessageQueue_push_blocked(t *testing.T) {
	msgQue := &messageQueue{}

	for i := 0; i < 10; i++ {
		msgQue.push(&testOperation{})
	}

	operation := &testOperation{}

	// Push operation and ensure that it's blocked
	msgQue.push(operation)

	if len(msgQue.ops) != 11 {
		t.Error("failed to push operation into queue")
	}

	if msgQue.blocked != nil {
		t.Error("expected message queue to be blocked")
	}

	// Remove operation from queue and ensure that it's unblocked
	msgQue.pull()

	if len(msgQue.ops) != 10 {
		t.Error("failed to pull operation from queue")
	}

	if msgQue.blocked != nil {
		t.Error("expected message queue to be unblocked")
	}

	if msgQue.empty() {
		t.Error("expected queue to not be empty")
	}
}

type testReceiver struct {
	LastTexture screen.Texture
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.LastTexture = t
}

type testOperation struct {
	updated bool
}

func (operation *testOperation) Do(t screen.Texture) bool {
	operation.updated = true
	return true
}

type mockScreen struct{}

func (m mockScreen) NewBuffer(image.Point) (screen.Buffer, error) {
	panic("implement me")
}

func (m mockScreen) NewTexture(image.Point) (screen.Texture, error) {
	return new(mockTexture), nil
}

func (m mockScreen) NewWindow(*screen.NewWindowOptions) (screen.Window, error) {
	panic("implement me")
}

type mockTexture struct {
	FillCnt int
}

func (m *mockTexture) Release() {}

func (m *mockTexture) Size() image.Point { return size }

func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: size}
}

func (m *mockTexture) Upload(image.Point, screen.Buffer, image.Rectangle) {
	panic("implement me")
}

func (m *mockTexture) Fill(image.Rectangle, color.Color, draw.Op) {
	m.FillCnt++
}
