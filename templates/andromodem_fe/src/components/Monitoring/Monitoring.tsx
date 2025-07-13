import { type FC } from "react"
import { toast } from "react-toastify"

import { monitoringApi } from "@/api/monitoringApi"
import useMonitoring from "@/hooks/useMonitoring"
import type { MonitoringFormData } from "@/hooks/useMonitoringForm"
import type { Device } from "@/types/device"
import type { MonitoringConfigPayload } from "@/types/monitoring"
import { showModal } from "@/utils/common"

import MonitoringConfigModal from "../Modal/MonitoringConfigModal"
import MonitoringConfigCard from "./MonitoringConfigCard"
import MonitoringLogCard from "./MonitoringLogCard"

const Monitoring: FC<{ device: Device }> = ({ device }) => {
  const {
    monitoringConfig,
    isLoadingMonitoringConfig,
    mutateMonitoringConfig,
    startMonitoring,
    stopMonitoring,
    isStarting,
    isStopping
  } = useMonitoring(device)

  const modalId = "monitoring-config-modal"

  const handleCreateConfig = () => {
    showModal(modalId)
  }

  const handleEdit = () => {
    if (monitoringConfig) {
      showModal(modalId)
    } else {
      console.warn('No monitoring config available for editing')
    }
  }

  const handleCancel = () => {
    const modal = document.getElementById(modalId) as HTMLDialogElement | null
    if (modal) {
      modal.close()
    }
  }

  const onSubmit = async (data: MonitoringFormData) => {
    try {
      const payload: MonitoringConfigPayload = {
        host: data.host,
        method: data.method,
        max_failures: data.max_failures,
        airplane_mode_delay: data.airplane_mode_delay,
        checking_interval: data.checking_interval
      }

      if (monitoringConfig) {
        await monitoringApi.updateDeviceMonitoringConfiguration(device.serial, payload)
        toast.success("Monitoring configuration updated successfully!")
      } else {
        await monitoringApi.createDeviceMonitoringConfiguration(device.serial, payload)
        toast.success("Monitoring configuration created successfully!")
      }

      await mutateMonitoringConfig()
      
      const modal = document.getElementById(modalId) as HTMLDialogElement | null
      if (modal) {
        modal.close()
      }
    } 
    // eslint-disable-next-line
    catch (error: any) {
      console.error('Failed to save monitoring config:', error)
      const errorMessage = error?.response?.data?.message || error?.message || 'Failed to save monitoring configuration'
      toast.error(errorMessage)
    }
  }

  if (isLoadingMonitoringConfig) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <div className="loading loading-spinner loading-lg"></div>
      </div>
    )
  }

  return (
    <div>
      <div className="col-span-1 md:col-span-2 grid grid-cols-1 gap-5 md:gap-10 justify-items-center mt-10">
        <MonitoringConfigCard
          monitoringConfig={monitoringConfig}
          onEdit={handleEdit}
          onCreate={handleCreateConfig}
          onStart={startMonitoring}
          onStop={stopMonitoring}
          isStarting={isStarting}
          isStopping={isStopping}
        />

        {monitoringConfig && (
          <MonitoringLogCard device={device} />
        )}
      </div>

      <MonitoringConfigModal
        modalId={modalId}
        monitoringConfig={monitoringConfig}
        onSubmit={onSubmit}
        onCancel={handleCancel}
      />
    </div>
  )
}

export default Monitoring
