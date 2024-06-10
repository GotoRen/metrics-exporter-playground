package pushmetric

import "github.com/prometheus/client_golang/prometheus"

// Collector は PushGateway に送信するカスタムメトリクスを管理します。
type Collector struct {
	asynCollectors []prometheus.Collector // 定常的にメトリクスを出力するコレクタです
	syncCollectors []prometheus.Collector // 任意のタイミングでメトリクスを出力するコレクタです
}

// NewCollector initializes a new Collector instance.
func NewCollector() *Collector {
	return &Collector{
		asynCollectors: make([]prometheus.Collector, 0),
		syncCollectors: make([]prometheus.Collector, 0),
	}
}

// RegisterAsyncMetrics はシーケンシャルに出力するカスタムメトリクスを登録します
// このメトリクスを pushInterval の値に基づいて定期的にメトリクスを出力します
func (c *Collector) RegisterAsyncMetrics(cms ...prometheus.Collector) *Collector {
	c.asynCollectors = append(c.asynCollectors, cms...)
	return c
}

// RegisterSyncMetrics は任意のタイミングで出力するメトリクスを登録します
// 任意のタイミングで Export 関数を呼び出すことでメトリクスを出力します
func (c *Collector) RegisterSyncMetrics(cms ...prometheus.Collector) *Collector {
	c.syncCollectors = append(c.syncCollectors, cms...)
	return c
}
