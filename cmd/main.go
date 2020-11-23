package main

import (
	// "encoding/json"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
	"github.com/prometheus/client_golang/tree/master/prometheus/promhttp"
)

var (
//verbMetrics *prometheus.GaugeVec
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
	Results map[string]float64 `json:"results"`
	Labels  map[string]string  `json:"labels"`
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
	labels := flag.String("labels", "hostname,env", "additioanal labels")
	//prefix := flag.String("prefix", "bash", "Prefix for metrics")
	debug := flag.Bool("debug", false, "Debug log level")
	flag.Parse()

	var labelsArr []string

	labelsArr = strings.Split(*labels, ",")
	labelsArr = append(labelsArr, "verb", "job")

	log.Println("Labels:")
	log.Println(labelsArr)

	//verbMetrics = prometheus.NewGaugeVec(
	//	prometheus.GaugeOpts{
	//		Name: fmt.Sprintf("%s", *prefix),
	//		Help: "bash exporter metrics",
	//	},
	// []string{"verb", "job"},
	//	labelsArr,
	//)
	//prometheus.MustRegister(verbMetrics)

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
		//verbMetrics.Reset()
		//if debug == true {
		//	ser, err := json.Marshal(o)
		//	if err != nil {
		//		log.Println(err)
		//	}
		//	log.Println(string(ser))
		//}

		for _, o := range oArr {

			for metric, value := range o.Schema.Results {
				for _, label := range labelsArr {
					if _, ok := o.Schema.Labels[label]; !ok {
						o.Schema.Labels[label] = ""
					}
				}
				o.Schema.Labels["verb"] = metric
				o.Schema.Labels["job"] = o.Job
				log.Println("verbMetrics")
				log.Println(o.Schema.Labels)
				log.Println(fmt.Sprintf("%f", float64(value)))
				//verbMetrics.With(prometheus.Labels(o.Schema.Labels)).Set(float64(value))
			}
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
