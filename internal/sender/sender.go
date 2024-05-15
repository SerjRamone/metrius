// Package sender ...
package sender

import (
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"

	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/SerjRamone/metrius/pkg/logger"
)

// metricsSender ...
type metricsSender struct {
	client  APIClient
	sURL    string
	hashKey string
	pubKey  []byte
	localIP string
}

// NewMetricsSender crates MetricsSender
func NewMetricsSender(sURL, hashKey string, pubKey []byte, serverType string) (*metricsSender, error) {
	if serverType != "grpc" && serverType != "http" {
		return nil, fmt.Errorf("server type %s not supported", serverType)
	}

	if serverType == "grpc" {
		gRPCClient, err := NewGRPCApiClient(sURL)
		if err != nil {
			logger.Error("can't create gRPC client", zap.Error(err))
			return nil, err
		}
		return &metricsSender{
			client: gRPCClient,
		}, nil
	}

	ip, err := getLocalIP()
	if err != nil {
		logger.Error("can't get local IP", zap.Error(err))
		return nil, err
	}
	logger.Info("client local IP", zap.String("IP", ip))

	httpClient := NewHTTPApiClient("http://"+sURL, ip)
	middlewares := make([]middleware, 0, 3)
	if hashKey != "" {
		middlewares = append(middlewares, hasher(hashKey))
	}
	if len(pubKey) != 0 {
		middlewares = append(middlewares, crypto(pubKey))
	}
	httpClient.client.Transport = chain(httpClient.client.Transport, middlewares...)

	sender := metricsSender{
		sURL:    "http://" + sURL,
		hashKey: hashKey,
		pubKey:  pubKey,
		client:  httpClient,
		localIP: ip,
	}

	return &sender, nil
}

// Send whole Collection
// func (sender *metricsSender) Send(collections []metrics.Collection) error {
// 	for _, c := range collections {
//
// 		for _, m := range c {
// 			r, err := sender.sendMetrics(m)
// 			if err != nil {
// 				return err
// 			}
//
// 			r.Body.Close()
// 		}
//
// 		time.Sleep(50 * time.Millisecond)
// 	}
//
// 	return nil
// }

func (sender *metricsSender) Send(collections []metrics.Collection) error {
	for _, c := range collections {

		for _, m := range c {
			err := sender.client.Do(m)
			if err != nil {
				return err
			}
		}

		time.Sleep(50 * time.Millisecond)
	}

	return nil
}

// Worker is a async sender
func (sender *metricsSender) Worker(doneCh chan struct{}, jobCh chan []metrics.Collection) {
	logger.Info("worker started")
	for {
		select {
		case collections := <-jobCh:
			logger.Info("worker recived new collections")
			err := sender.client.DoBatch(collections)
			if err != nil {
				logger.Error("async SendBatch error", zap.Error(err))
			}
		case <-doneCh:
			logger.Info("worker recived done signal")
			collections := <-jobCh
			if len(collections) > 0 {
				err := sender.client.DoBatch(collections)
				if err != nil {
					logger.Error("async SendBatch error", zap.Error(err))
				}
			}
			return
		}
	}
}

// SendBatch sends metrics in batches
// func (sender *metricsSender) SendBatch(collections []metrics.Collection) error {
// 	batch := make([]metrics.Metrics, 0, 200)
// 	// collect batch of metrics.Metrics
// 	for _, c := range collections {
// 		for _, m := range c {
// 			item := metrics.Metrics{
// 				ID:    m.Name,
// 				MType: m.Type,
// 			}
//
// 			switch m.Type {
// 			case "gauge":
// 				value := m.Value
// 				item.Value = &value
// 			case "counter":
// 				tmp := int64(m.Value)
// 				item.Delta = &tmp
// 			}
//
// 			batch = append(batch, item)
// 		}
// 	}
//
// 	if len(batch) > 0 {
// 		url := fmt.Sprintf("%s/updates/", sender.sURL)
// 		data, err := json.Marshal(batch)
// 		if err != nil {
// 			logger.Error("metrics encode error", zap.Error(err))
// 			return err
// 		}
//
// 		var b bytes.Buffer
// 		gw := gzip.NewWriter(&b)
// 		if _, err = gw.Write(data); err != nil {
// 			logger.Error("gzipped write data error", zap.Error(err))
// 			return err
// 		}
// 		if err = gw.Close(); err != nil {
// 			return err
// 		}
//
// 		req, err := http.NewRequest("POST", url, &b)
// 		if err != nil {
// 			logger.Error("request object creation error", zap.Error(err))
// 			return err
// 		}
// 		req.Header.Set("Content-Type", "application/json")
// 		req.Header.Set("Content-Encoding", "gzip")
// 		req.Header.Set("X-Real-IP", sender.localIP)
//
// 		r, err := sender.client.Do(req)
// 		if err != nil {
//
// 			logger.Error("first request error", zap.Error(err))
// 			var netError net.Error
// 			if errors.As(err, &netError) {
// 				err = retry.WithBackoff(func() error {
// 					retrResp, retrErr := sender.client.Do(req)
// 					if retrErr != nil {
// 						return retrErr
// 					}
// 					_, retrErr = io.Copy(io.Discard, retrResp.Body)
// 					if retrErr != nil {
// 						logger.Error("io.Copy error", zap.Error(retrErr))
// 					}
// 					defer retrResp.Body.Close()
// 					return retrErr
// 				}, 3)
// 				return err
// 			}
// 			logger.Error("send metrics in batch error", zap.Error(err))
// 			return err
// 		}
// 		defer r.Body.Close()
//
// 		_, err = io.Copy(io.Discard, r.Body)
// 		if err != nil {
// 			logger.Error("io.Copy error", zap.Error(err))
// 			return err
// 		}
// 	}
//
// 	return nil
// }

// single request
// func (sender *metricsSender) sendMetrics(m metrics.CollectionItem) (*http.Response, error) {
// 	url := fmt.Sprintf("%s/update/", sender.sURL)
//
// 	item := metrics.Metrics{
// 		ID:    m.Name,
// 		MType: m.Type,
// 	}
//
// 	switch m.Type {
// 	case "gauge":
// 		item.Value = &m.Value
// 	case "counter":
// 		tmp := int64(m.Value)
// 		item.Delta = &tmp
// 	}
//
// 	data, err := json.Marshal(item)
// 	if err != nil {
// 		logger.Error("metrics encode error", zap.Error(err))
// 	}
//
// 	var b bytes.Buffer
// 	gw := gzip.NewWriter(&b)
// 	if _, err = gw.Write(data); err != nil {
// 		return nil, err
// 	}
// 	if err = gw.Close(); err != nil {
// 		return nil, err
// 	}
// 	req, err := http.NewRequest("POST", url, &b)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Content-Encoding", "gzip")
// 	req.Header.Set("X-Real-IP", sender.localIP)
//
// 	r, err := sender.client.Do(req)
// 	if err != nil {
// 		var netError net.Error
// 		if errors.As(err, &netError) {
// 			err = retry.WithBackoff(func() error {
// 				retrResp, retrErr := sender.client.Do(req)
// 				if retrErr != nil {
// 					_, retrErr = io.Copy(io.Discard, retrResp.Body)
// 					if retrErr != nil {
// 						return retrErr
// 					}
// 				}
// 				defer retrResp.Body.Close()
// 				return retrErr
// 			}, 3)
// 			return r, err
// 		}
// 		logger.Error("send metrics error", zap.Error(err))
// 		return nil, err
// 	}
// 	defer r.Body.Close()
//
// 	_, err = io.Copy(io.Discard, r.Body)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return r, nil
// }

// getLocalIP return local IPv4 client address
func getLocalIP() (string, error) {
	// get IP addresses
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("gettin IP error")
	}

	var ipAddress string
	for _, address := range addrs {
		ipnet, ok := address.(*net.IPNet)
		// checking if the address is IPv4 and not loopback
		if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ipAddress = ipnet.IP.String()
			break
		}
	}

	if ipAddress == "" {
		return "", fmt.Errorf("local IPv4 not found")
	}

	return ipAddress, nil
}
