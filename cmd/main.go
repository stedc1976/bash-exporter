package main

import (
	// "encoding/json"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	bashMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "bash",
		Help: "bash exporter metrics",
	}, []string{"pod_name", "namespace", "container_name"},
	)
)

// Type Params stores parameters.
type Params struct {
	Path  *string
	UseWg bool
	Wg    *sync.WaitGroup
}

type Output struct {
	Schema Schema `json:""`
	Job    string `json:""`
}

type Schema struct {
	Results map[string]int64  `json:"results"`
	Labels  map[string]string `json:"labels"`
}

func (o *Output) RunJob(p *Params) {
	if p.UseWg {
		defer p.Wg.Done()
	}
	o.RunExec(p.Path)
}

func (o *Output) RunExec(path *string) {

	out, err := exec.Command(*path).Output()
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(out, &o.Schema)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	addr := flag.String("web.listen-address", ":9300", "Address on which to expose metrics")
	interval := flag.Int("interval", 5, "Interval for metrics collection in seconds")
	path := flag.String("path", "./scripts", "path to directory with bash scripts")
	debug := flag.Bool("debug", true, "Debug log level")
	flag.Parse()

	labelsArr := []string{"pod_name", "namespace", "container_name"}

	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(bashMetric)

	files, err := ioutil.ReadDir(*path)
	if err != nil {
		log.Fatal(err)
	}

	var names []string
	for _, f := range files {
		if f.Name()[0:1] != "." {
			names = append(names, f.Name())
		}
	}

	log.Println("Bash scripts found:")
	log.Println(names)

	http.Handle("/metrics", promhttp.Handler())
	go Run(int(*interval), *path, names, labelsArr, *debug)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func Run(interval int, path string, names []string, labelsArr []string, debug bool) {
	for {
		var wg sync.WaitGroup
		oArr := []*Output{}
		wg.Add(len(names))
		for _, name := range names {
			o := Output{}
			o.Job = strings.Split(name, ".")[0]
			oArr = append(oArr, &o)
			thisPath := path + "/" + name
			p := Params{UseWg: true, Wg: &wg, Path: &thisPath}
			go o.RunJob(&p)
		}
		wg.Wait()
		bashMetric.Reset()

		for _, o := range oArr {
			for metric, value := range o.Schema.Results {
				for _, label := range labelsArr {
					if _, ok := o.Schema.Labels[label]; !ok {
						o.Schema.Labels[label] = ""
					}
				}
				o.Schema.Labels["verb"] = metric
				o.Schema.Labels["job"] = o.Job

				if debug == true {
					log.Println(o.Schema.Labels)
					log.Println(int64(value))
				}

				bashMetric.With(prometheus.Labels(o.Schema.Labels)).Set(float64(int64(value)))
			}
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
