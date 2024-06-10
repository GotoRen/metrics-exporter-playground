package pushmetric

// // TestUpdateDefaultMetric tests the updateDefaultMetric function.
// func TestUpdateDefaultMetric(t *testing.T) {
// 	applicationName := "test-app"
// 	instance := "test-instance"
// 	labels := []string{applicationName, instance}

// 	initialCPUUtilization, err := getCPUUtilization()
// 	assert.NoError(t, err, "Error should be nil when getting initial CPU utilization")
// 	assert.GreaterOrEqual(t, initialCPUUtilization, 0.0, "Initial CPU Utilization should be non-negative")
// 	assert.LessOrEqual(t, initialCPUUtilization, 100.0, "Initial CPU Utilization should be less than or equal to 100")

// 	initialMemoryUtilization, err := getMemoryUtilization()
// 	assert.NoError(t, err, "Error should be nil when getting initial memory utilization")
// 	assert.GreaterOrEqual(t, initialMemoryUtilization, 0.0, "Initial Memory Utilization should be non-negative")
// 	assert.LessOrEqual(t, initialMemoryUtilization, 100.0, "Initial Memory Utilization should be less than or equal to 100")

// 	// Update metrics.
// 	err = updateDefaultMetric(labels...)
// 	assert.NoError(t, err, "Error should be nil when updating metrics")

// 	// Verify that the metrics have been updated with expected value range.
// 	actualCPUUtilizationValue := testutil.ToFloat64(cpuUtilizationMetric.WithLabelValues(applicationName, instance))
// 	assert.GreaterOrEqual(t, actualCPUUtilizationValue, 0.0, "Actual CPU Utilization should be non-negative")
// 	assert.LessOrEqual(t, actualCPUUtilizationValue, 100.0, "Actual CPU Utilization should be less than or equal to 100")

// 	actualMemoryUtilizationValue := testutil.ToFloat64(memoryUtilizationMetric.WithLabelValues(applicationName, instance))
// 	assert.GreaterOrEqual(t, actualMemoryUtilizationValue, 0.0, "Actual Memory Utilization should be non-negative")
// 	assert.LessOrEqual(t, actualMemoryUtilizationValue, 100.0, "Actual Memory Utilization should be less than or equal to 100")

// 	// Verify that the push count has been incremented.
// 	actualPushCount := testutil.ToFloat64(pushCountMetric.WithLabelValues(applicationName, instance))
// 	assert.Equal(t, 1.0, actualPushCount, "Push count metric should be incremented to 1")
// }

// // TestGetCPUUtilization tests the getCPUUtilization function.
// func TestGetCpuUtilization(t *testing.T) {
// 	cpuUtilizationValue, err := getCPUUtilization()
// 	assert.NoError(t, err, "Error should be nil when getting CPU utilization")
// 	assert.GreaterOrEqual(t, cpuUtilizationValue, 0.0, "CPU Utilization should be a non-negative value")
// 	assert.LessOrEqual(t, cpuUtilizationValue, 100.0, "CPU Utilization should not exceed 100")
// }

// // TestGetMemoryUtilization tests the getMemoryUtilization function.
// func TestGetMemoryUtilization(t *testing.T) {
// 	memoryUtilizationValue, err := getMemoryUtilization()
// 	assert.NoError(t, err, "Error should be nil when getting memory utilization")
// 	assert.GreaterOrEqual(t, memoryUtilizationValue, 0.0, "Memory Utilization should be a non-negative value")
// 	assert.LessOrEqual(t, memoryUtilizationValue, 100.0, "Memory Utilization should not exceed 100")
// }
