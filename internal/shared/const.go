package shared

const (
	GaugeMetricType   = "gauge"
	CounterMetricType = "counter"
)

// GaugeMetricName определяет набор возможных метрик-значений по умолчанию.
type GaugeMetricName string

const (
	Alloc         = GaugeMetricName("Alloc")
	BuckHashSys   = GaugeMetricName("BuckHashSys")
	FreeMemory    = GaugeMetricName("FreeMemory")
	Frees         = GaugeMetricName("Frees")
	GcCPUFraction = GaugeMetricName("GCCPUFraction")
	GcSys         = GaugeMetricName("GCSys")
	HeapAlloc     = GaugeMetricName("HeapAlloc")
	HeapIdle      = GaugeMetricName("HeapIdle")
	HeapInuse     = GaugeMetricName("HeapInuse")
	HeapObjects   = GaugeMetricName("HeapObjects")
	HeapReleased  = GaugeMetricName("HeapReleased")
	HeapSys       = GaugeMetricName("HeapSys")
	LastGc        = GaugeMetricName("LastGC")
	Lookups       = GaugeMetricName("Lookups")
	MemCacheInuse = GaugeMetricName("MCacheInuse")
	MemCacheSys   = GaugeMetricName("MCacheSys")
	MemSpanInuse  = GaugeMetricName("MSpanInuse")
	MemSpanSys    = GaugeMetricName("MSpanSys")
	MemAllocs     = GaugeMetricName("Mallocs")
	NextGc        = GaugeMetricName("NextGC")
	NumForcedGc   = GaugeMetricName("NumForcedGC")
	NumGc         = GaugeMetricName("NumGC")
	OtherSys      = GaugeMetricName("OtherSys")
	PauseTotalNs  = GaugeMetricName("PauseTotalNs")
	RandomValue   = GaugeMetricName("RandomValue")
	StackInuse    = GaugeMetricName("StackInuse")
	StackSys      = GaugeMetricName("StackSys")
	Sys           = GaugeMetricName("Sys")
	TotalAlloc    = GaugeMetricName("TotalAlloc")
	TotalMemory   = GaugeMetricName("TotalMemory")
)

// CounterMetricName определяет набор возможных метрик-счётчиков по умолчанию.
type CounterMetricName string

const (
	PollCount = CounterMetricName("PollCount")
)
