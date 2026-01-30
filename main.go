package main

import (
	"context"
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//go:embed index.html
var templateFS embed.FS

type ContainersInfo struct {
	Name        string
	Id          string
	Status      string
	State       container.ContainerState
	CPUUsage    float64
	MemoryUsage uint64
	Port        uint16
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	port := os.Getenv("port")
	title := os.Getenv("title")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	templ := template.Must(template.New("").ParseFS(templateFS, "index.html"))
	r.SetHTMLTemplate(templ)

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": title,
			"data":  getContainers(),
		})
	})

	r.POST("/toggle/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		err := toggleContainer(id)
		log.Printf("Error toggling for %s:%e", id, err)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.Status(http.StatusOK)
	})

	log.Println("Server Running!")
	if err := r.Run(port); err != nil {
		panic(err)
	}
}

// DOCKER STUFF -------------------------------

func getContainers() []ContainersInfo {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true, Filters: filters.NewArgs(filters.Arg("name", "mc-"))})
	if err != nil {
		panic(err)
	}

	var containersInfo []ContainersInfo

	for _, ctr := range containers {
		resp, err := apiClient.ContainerStats(context.Background(), ctr.ID, false)
		if err != nil {
			log.Printf("Error getting stats for container %s: %v", ctr.ID, err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading stats for container %s: %v", ctr.ID, err)
			continue
		}

		var stats container.StatsResponse
		if err := json.Unmarshal(body, &stats); err != nil {
			log.Printf("Error decoding stats for container %s: %v", ctr.ID, err)
			continue
		}

		// 1. Calculate the deltas
		cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
		systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)

		// 2. Get the number of CPUs (OnlineCPUs is more reliable)
		numberCPUs := float64(stats.CPUStats.OnlineCPUs)
		if numberCPUs == 0 {
			// Fallback if OnlineCPUs isn't reported
			numberCPUs = float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
		}

		cpuPercent := 0.0
		if systemDelta > 0.0 && cpuDelta > 0.0 {
			// 3. The Formula
			cpuPercent = (cpuDelta / systemDelta) * numberCPUs * 100.0
		}

		var port uint16 = 0
		if len(ctr.Ports) > 0 {
			port = ctr.Ports[0].PublicPort
		}

		containersInfo = append(containersInfo, ContainersInfo{
			Name:        ctr.Names[0],
			Id:          ctr.ID,
			Status:      ctr.Status,
			State:       ctr.State,
			MemoryUsage: stats.MemoryStats.Usage / 1000000,
			CPUUsage:    cpuPercent,
			Port:        port,
		})
	}

	return containersInfo
}

func toggleContainer(id string) error {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	inspect, err := apiClient.ContainerInspect(context.Background(), id)
	if err != nil {
		return err
	}
	if inspect.State.Running {
		return apiClient.ContainerStop(context.Background(), id, container.StopOptions{})
	} else {
		return apiClient.ContainerStart(context.Background(), id, container.StartOptions{})
	}
}
