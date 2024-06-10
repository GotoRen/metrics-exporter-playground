package pushmetric

// // TestNew verifies that a new Exporter instance is initialized correctly.
// func TestNew(t *testing.T) {
// 	jobName := "test-job"
// 	applicationName := "test-app"
// 	pushInterval := 1 * time.Second
// 	endPoint := "http://example.com"
// 	collector := NewCollector()

// 	exporter := New(jobName, applicationName, pushInterval, endPoint, collector)

// 	assert.Equal(t, jobName, exporter.jobName)
// 	assert.Equal(t, applicationName, exporter.applicationName)
// 	assert.Equal(t, pushInterval, exporter.pushInterval)
// 	assert.Equal(t, endPoint, exporter.endPoint)
// 	assert.Equal(t, collector, exporter.collector)
// 	assert.NotNil(t, exporter.client)
// }

// // TestWithClient verifies that the HTTP client is correctly set in the Exporter.
// func TestWithClient(t *testing.T) {
// 	exporter := New("test-job", "test-app", 1*time.Second, "http://example.com", NewCollector())

// 	customClient := &http.Client{
// 		Timeout:   5 * time.Second,
// 		Transport: http.DefaultTransport.(*http.Transport).Clone(),
// 	}

// 	exporter.WithClient(customClient)

// 	assert.Equal(t, customClient, exporter.client)
// }

// // TestExport verifies that the Export method correctly pushes metrics to the Pushgateway.
// func TestExport(t *testing.T) {
// 	// Setup a mock Pushgateway server.
// 	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 	}))
// 	defer mockServer.Close()

// 	ctx := context.Background()

// 	collector := NewCollector()

// 	exporter := New("test-job", "test-app", 1*time.Second, mockServer.URL, collector)

// 	err := exporter.Export(ctx)
// 	assert.NoError(t, err, "Export should succeed")
// }

// // TestRoutineSequentialExporter verifies the continuous metrics collection and pushing.
// func TestRoutineSequentialExporter(t *testing.T) {
// 	// Setup a mock Pushgateway server
// 	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		time.Sleep(200 * time.Millisecond) // Simulate delay in response
// 		w.WriteHeader(http.StatusOK)
// 	}))
// 	defer mockServer.Close()

// 	collector := NewCollector()
// 	exporter := New("test-job", "test-app", 100*time.Millisecond, mockServer.URL, collector)

// 	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
// 	defer cancel()

// 	errCh := make(chan error, 1)
// 	go func() {
// 		errCh <- exporter.RoutineSequentialExporter(ctx)
// 	}()

// 	time.Sleep(400 * time.Millisecond) // Wait for CronJob execution time

// 	// Verify that the function exits with context deadline exceeded error.
// 	err := <-errCh
// 	if assert.Error(t, err, "RoutineSequentialExporter should return an error") {
// 		assert.True(t, errors.Is(err, context.DeadlineExceeded), "RoutineSequentialExporter should exit with context deadline exceeded error: %v", err)
// 	}
// }
