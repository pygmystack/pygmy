package docker

import (
	"context"
	"fmt"
	"strconv"

	"github.com/docker/docker/client"
)

// SetField will set a pygmy label to be equal to the string equal of
// an interface{}, even if it already exists. It should not matter if
// this container is running or not.
func (Service *Service) SetField(ctx context.Context, cli *client.Client, name string, value interface{}) error {
	if _, ok := Service.Config.Labels["pygmy."+fmt.Sprint(name)]; !ok {
		//
	} else {
		old, _ := Service.GetFieldString(ctx, cli, name)
		Service.Config.Labels["pygmy."+name] = fmt.Sprint(value)
		new, _ := Service.GetFieldString(ctx, cli, name)

		if old == new {
			return fmt.Errorf("tag was not set")
		}
	}

	return nil
}

// GetFieldString will get and return a tag on the service using the pygmy
// convention ("pygmy.*") and return it as a string.
func (Service *Service) GetFieldString(ctx context.Context, cli *client.Client, field string) (string, error) {

	f := fmt.Sprintf("pygmy.%v", field)

	if labels, running := Service.Labels(ctx, cli); running == nil {
		if val, ok := labels[f]; ok {
			return val, nil
		}
	}

	if val, ok := Service.Config.Labels[f]; ok {
		return val, nil
	}

	return "", fmt.Errorf("could not find field 'pygmy.%v' on service using image %v?", field, Service.Config.Image)
}

// GetFieldInt will get and return a tag on the service using the pygmy
// convention ("pygmy.*") and return it as an int.
func (Service *Service) GetFieldInt(ctx context.Context, cli *client.Client, field string) (int, error) {

	f := fmt.Sprintf("pygmy.%v", field)

	if labels, running := Service.Labels(ctx, cli); running == nil {
		if val, ok := labels[f]; ok {
			i, e := strconv.ParseInt(val, 10, 10)
			if e != nil {
				return 0, e
			}
			return int(i), nil
		}
	}

	if val, ok := Service.Config.Labels[f]; ok {
		i, e := strconv.ParseInt(val, 10, 10)
		if e != nil {
			return 0, e
		}
		return int(i), nil
	}

	return 0, fmt.Errorf("could not find field 'pygmy.%v' on service using image %v?", field, Service.Config.Image)
}

// GetFieldBool will get and return a tag on the service using the pygmy
// convention ("pygmy.*") and return it as a bool.
func (Service *Service) GetFieldBool(ctx context.Context, cli *client.Client, field string) (bool, error) {

	f := fmt.Sprintf("pygmy.%v", field)

	if labels, running := Service.Labels(ctx, cli); running == nil {
		if Service.Config.Labels[f] == labels[f] {
			if val, ok := labels[f]; ok {
				if val == "true" {
					return true, nil
				} else if val == "false" {
					return false, nil
				}
			}
		}
	}

	if val, ok := Service.Config.Labels[f]; ok {
		if val == "true" || val == "1" {
			return true, nil
		} else if val == "false" || val == "0" {
			return false, nil
		}
	}

	return false, fmt.Errorf("could not find field 'pygmy.%v' on service using image %v?", field, Service.Config.Image)
}
