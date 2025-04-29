package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	UsersCreationCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_create_total",
			Help: "Number of created users",
		},
		[]string{"method", "status"},
	)

	UsersCreationDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "user_create_duration_seconds",
			Help: "Duration of creating users",
		},
		[]string{"method"},
	)

	UsersCreationErrorsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_create_errors_total",
			Help: "Number of failed users creation",
		},
		[]string{"method", "reason"},
	)

	UsersLoginTotalCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_login_total",
			Help: "Number of total login attempts",
		},
		[]string{"method", "status"},
	)

	UsersLoginDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "user_login_duration_seconds",
			Help: "Duration of login attempts",
		},
		[]string{"method"},
	)

	UserLogoutTotalCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "user_logout_total",
			Help: "Number of total logout attempts",
		})
)

var (
	ProfileDeletionTotalCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "profile_delete_total",
			Help: "Number of deleted profiles",
		},
	)

	ProfileChangePasswordTotalCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "profile_change_password_total",
			Help: "Number of changed profiles passwords",
		},
	)
)

var (
	AdminActionsTotalCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "admin_actions_total",
			Help: "Number of admin actions",
		},
		[]string{"method"},
	)

	AdminForbittenActionsTotalCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "admin_forbitten_actions_total",
			Help: "Number of forbitten admin actions",
		},
		[]string{"method", "reason"})
)

func Initialize() {
	prometheus.MustRegister(UsersCreationCounter)
	prometheus.MustRegister(UsersCreationDurationHistogram)
	prometheus.MustRegister(UsersCreationErrorsCounter)
	prometheus.MustRegister(UsersLoginTotalCounter)
	prometheus.MustRegister(UsersLoginDurationHistogram)
	prometheus.MustRegister(UserLogoutTotalCounter)

	prometheus.MustRegister(ProfileDeletionTotalCounter)
	prometheus.MustRegister(ProfileChangePasswordTotalCounter)

	prometheus.MustRegister(AdminActionsTotalCounter)
	prometheus.MustRegister(AdminForbittenActionsTotalCounter)
}
