package main

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	sway "github.com/joshuarubin/go-sway"
)

type handler struct {
	sway.EventHandler
	client sway.Client
}

type workspaceName struct {
	number    int
	shortName *string
	icon      *string
}

func getFocuesdWorkspace(workspaces *[]sway.Workspace) *sway.Workspace {
	for _, workspace := range *workspaces {
		if workspace.Focused {
			return &workspace
		}
	}

	return nil
}

func (h handler) renameWorkspace(ctx context.Context, currentName, newName string) ([]sway.RunCommandReply, error) {
	return h.client.RunCommand(ctx, fmt.Sprintf("rename workspace %s to %s", currentName, newName))
}

func parseWorkspaceName(name string) (workspaceName, error) {
	regex := regexp.MustCompile(`(?P<num>[0-9]+):?(?P<shortname>\w+)? ?(?P<icons>.+)?`)

	matches := regex.FindAllStringSubmatch(name, -1)
	if matches == nil {
		return workspaceName{}, errors.New("No matches in workspace name")
	}

	shortName := matches[0][0]
	icon := matches[0][2]

	number, err := strconv.Atoi(matches[0][1])
	if err != nil {
		return workspaceName{}, err
	}

	return workspaceName{
		number:    number,
		shortName: &shortName,
		icon:      &icon,
	}, nil
}

func (h handler) Window(ctx context.Context, e sway.WindowEvent) {
	workspaces, err := h.client.GetWorkspaces(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	focusedWorkspace := getFocuesdWorkspace(&workspaces)
	if focusedWorkspace == nil {
		return
	}

	workspaceName, err := parseWorkspaceName(focusedWorkspace.Name)
	if err != nil {
		fmt.Println(err)
		return
	}

	name := strings.Split(e.Container.Name, " ")[0]
	_, err = h.renameWorkspace(ctx, focusedWorkspace.Name, fmt.Sprintf("%d: %s", workspaceName.number, name))

	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	ctx := context.TODO()

	client, err := sway.New(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	eventHandler := handler{
		EventHandler: sway.NoOpEventHandler(),
		client:       client,
	}

	err = sway.Subscribe(ctx, eventHandler, sway.EventTypeWindow)
	if err != nil {
		fmt.Println(err)
	}
}
