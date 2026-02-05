package messaging

const (
	// Exchanges
	CrawlCommandExchange = "crawl.command.x"
	CrawlEventExchange   = "crawl.event.x"

	// Queues
	CrawlerExecuteQueue  = "crawler.crawl.execute.q"
	SchedulerResultQueue = "scheduler.crawl.result.q"

	// Routing keys
	CrawlCommandExecuteKey = "crawl.command.execute"
	CrawlEventSucceededKey = "crawl.event.succeeded"
	CrawlEventFailedKey    = "crawl.event.failed"
)
