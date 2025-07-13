package rest

import "net/http"

type IMonitoringHandler interface {
	CreateMonitoring(w http.ResponseWriter, r *http.Request)
	StartMonitoring(w http.ResponseWriter, r *http.Request)
	StopMonitoring(w http.ResponseWriter, r *http.Request)
	DeleteMonitoring(w http.ResponseWriter, r *http.Request)
	ClearMonitoringLogs(w http.ResponseWriter, r *http.Request)
	GetMonitoringStatus(w http.ResponseWriter, r *http.Request)
	GetMonitoringConfig(w http.ResponseWriter, r *http.Request)
	UpdateMonitoringConfig(w http.ResponseWriter, r *http.Request)
	GetAllMonitoringTasks(w http.ResponseWriter, r *http.Request)
	GetMonitoringLogs(w http.ResponseWriter, r *http.Request)
}
