import { MonitoringMethod, type MonitoringMethodValue } from "@/types/monitoring"

export const monitoringMethods = [
  { value: MonitoringMethod.HTTP, label: 'HTTP', description: 'Standard HTTP ping' },
  { value: MonitoringMethod.HTTPS, label: 'HTTPS', description: 'Secure HTTP ping' },
  // { value: MonitoringMethod.WS, label: 'WebSocket', description: 'WebSocket connection test' },
  { value: MonitoringMethod.PING, label: 'Ping', description: 'ICMP ping' },
  { value: MonitoringMethod.PING_BY_DEVICE, label: 'Ping by device', description: 'ICMP Ping using adb command directly on device' }
]

export const getProtocolFromMethod = (method: MonitoringMethodValue): string => {
  switch (method) {
    case MonitoringMethod.HTTP:
      return 'http://'
    case MonitoringMethod.HTTPS:
      return 'https://'
    // case MonitoringMethod.WS:
    //   return 'ws://';
    case MonitoringMethod.PING:
      return 'ICMP'
    case MonitoringMethod.PING_BY_DEVICE:
      return 'ICMP'
    default:
      return ''
  }
}