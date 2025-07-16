package command

const (
	// Minimum Android version where showing APN info is supported without root access.
	// Below Android 11 (R), accessing APN settings may still be possible but typically requires root privileges.
	MinimumAndroidShowApn = 11

	// Minimum Android version where accessing SMS/MMS messages via shell commands is supported without root.
	// On versions below Android 10 (Q), message access may still work, but usually needs root access.
	MinimumAndroidShowMessages = 10

	// Minimum Android version where toggling airplane mode is supported via `adb shell cmd connectivity airplane-mode` without requiring root access. On Android 8 (Oreo) and below, airplane mode can only be toggled by sendinga broadcast intent, which requires root permission.
	MinimumAndroidToggleAirplaneMode = 9

	// Minimum Android version where toggling mobile data via `svc data` is supported.
	// While some community references and AOSP commits suggest that the `svc` command
	// was introduced around Android 4.0 (Ice Cream Sandwich), in practice, reliable support
	// Assumption: We assume `svc` is present and functional starting in Android 5,
	// especially for mobile data toggling. Behavior may vary across OEMs and ROM versions
	MinimumAndroidToggleMobileData = 5

	// Minimum Android version where signal strength data from `dumpsys telephony.registry` can be reliably parsed. On Android versions below 10 (Q), the output format is inconsistent or incomplete, making it difficult to extract accurate signal strength values programmatically
	MinimumAndroidGetSignalStrength = 10
)
