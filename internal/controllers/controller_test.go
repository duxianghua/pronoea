package controllers

import (
	"context"
	"fmt"
	"testing"
	"time"

	probev1 "github.com/duxianghua/pronoea/internal/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestXxx(t *testing.T) {
	StartMgr()
	probeList := probev1.ProbeList{}
	time.Sleep(time.Second * 3)
	Probe.List(context.TODO(), &probeList, &client.ListOptions{})
	fmt.Println(probeList)
}
