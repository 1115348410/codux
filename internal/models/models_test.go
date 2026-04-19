package models

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewProject(t *testing.T) {
	project := NewProject("test", "/path/to/test")
	
	if project.Name != "test" {
		t.Errorf("Expected name 'test', got '%s'", project.Name)
	}
	
	if project.Path != "/path/to/test" {
		t.Errorf("Expected path '/path/to/test', got '%s'", project.Path)
	}
	
	if project.Shell != "bash" {
		t.Errorf("Expected shell 'bash', got '%s'", project.Shell)
	}
	
	if project.BadgeColorHex != "#007AFF" {
		t.Errorf("Expected color '#007AFF', got '%s'", project.BadgeColorHex)
	}
}

func TestNewProjectWorkspace(t *testing.T) {
	u := UUID{UUID: uuid.New()}
	workspace := NewProjectWorkspace(u)
	
	if workspace.ProjectID.String() != u.String() {
		t.Errorf("Expected projectID to match")
	}
	
	if workspace.Sessions == nil {
		t.Errorf("Expected sessions to be initialized")
	}
}

func TestDefaultAppSettings(t *testing.T) {
	settings := DefaultAppSettings()
	
	if settings.Theme != "system" {
		t.Errorf("Expected theme 'system', got '%s'", settings.Theme)
	}
	
	if settings.TerminalFontSize != 14 {
		t.Errorf("Expected fontSize 14, got %f", settings.TerminalFontSize)
	}
}

func TestSplitLayout(t *testing.T) {
	layout := NewSingleLayout("session1")
	
	if layout.Type != LayoutSingle {
		t.Errorf("Expected type LayoutSingle")
	}
	
	if layout.SessionID != "session1" {
		t.Errorf("Expected sessionID 'session1'")
	}
	
	ids := layout.GetSessionIDs()
	if len(ids) != 1 || ids[0] != "session1" {
		t.Errorf("Expected 1 sessionID")
	}
}

func TestSplitLayoutHSplit(t *testing.T) {
	left := NewSingleLayout("session1")
	right := NewSingleLayout("session2")
	layout := NewHSplitLayout(left, right, 0.5)
	
	if layout.Type != LayoutHSplit {
		t.Errorf("Expected type LayoutHSplit")
	}
	
	ids := layout.GetSessionIDs()
	if len(ids) != 2 {
		t.Errorf("Expected 2 sessionIDs")
	}
}

func TestSplitLayoutVSplit(t *testing.T) {
	top := NewSingleLayout("session1")
	bottom := NewSingleLayout("session2")
	layout := NewVSplitLayout(top, bottom, 0.5)
	
	if layout.Type != LayoutVSplit {
		t.Errorf("Expected type LayoutVSplit")
	}
	
	ids := layout.GetSessionIDs()
	if len(ids) != 2 {
		t.Errorf("Expected 2 sessionIDs")
	}
}
