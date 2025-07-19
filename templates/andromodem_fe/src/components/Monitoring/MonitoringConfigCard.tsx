import type { FC } from "react"
import { FaCircleStop } from "react-icons/fa6"
import { IoIosPlayCircle } from "react-icons/io"

import { useDevicesStore } from "@/stores/devicesStore"
import type { MonitoringConfig } from "@/types/monitoring"
import { isFeatureAvailable } from "@/utils/common"
import { monitoringMethods } from "@/utils/monitoringUtils"

interface MonitoringConfigCardProps {
  monitoringConfig: MonitoringConfig | null;
  onEdit: () => void;
  onCreate: () => void;
  onStart: () => void;
  onStop: () => void;
  isStarting: boolean;
  isStopping: boolean;
}

const MonitoringConfigCard: FC<MonitoringConfigCardProps> = ({
  monitoringConfig,
  onEdit,
  onCreate,
  onStart,
  onStop,
  isStarting,
  isStopping
}) => {
  const deviceFeatureAvailabilities = useDevicesStore((state) => state.deviceFeatureAvailabilities)
  const allowed = isFeatureAvailable(deviceFeatureAvailabilities, "can_change_airplane_mode_status")
  return (
    <div className="card w-full bg-base-100 shadow-sm">
      <div className="card-body">
        <div className="flex justify-between items-center mb-4">
          <h2 className="card-title">Configuration</h2>
          {monitoringConfig ? (
            <button
              onClick={onEdit}
              className="btn btn-primary btn-sm"
            >
              Edit Configuration
            </button>
          ) : <></>}
        </div>

        {monitoringConfig ? (
          <div className="space-y-4">
            {
              !allowed && !monitoringConfig.is_active ?
                <div role="alert" className="alert alert-warning">
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 shrink-0 stroke-current" fill="none" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <span>You can't start monitoring config. please check device feature availability.</span>
                </div>
                : <></>
            }
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <span className="text-sm font-medium">Status:</span>
                {
                  monitoringConfig.is_active ? (
                    <div className="badge badge-success">
                      <svg className="size-[1em]" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><g fill="currentColor" strokeLinejoin="miter" strokeLinecap="butt"><circle cx="12" cy="12" r="10" fill="none" stroke="currentColor" strokeLinecap="square" stroke-miterlimit="10" strokeWidth="2"></circle><polyline points="7 13 10 16 17 8" fill="none" stroke="currentColor" strokeLinecap="square" stroke-miterlimit="10" strokeWidth="2"></polyline></g></svg>
                      Running
                    </div>
                  ) : (
                    <div className="badge badge-error">
                      <svg className="size-[1em]" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><g fill="currentColor" strokeLinejoin="miter" strokeLinecap="butt"><circle cx="12" cy="12" r="10" fill="none" stroke="currentColor" strokeLinecap="square" stroke-miterlimit="10" strokeWidth="2"></circle><line x1="12" y1="8" x2="12" y2="12" fill="none" stroke="currentColor" strokeLinecap="square" stroke-miterlimit="10" strokeWidth="2"></line><line x1="12" y1="16" x2="12" y2="16" fill="none" stroke="currentColor" strokeLinecap="square" stroke-miterlimit="10" strokeWidth="2"></line></g></svg>
                      Stopped
                    </div>
                  )
                }
              </div>
            </div>
            <div>
              <label className="label">
                <span className="label-text font-semibold">Monitoring Method</span>
              </label>
              <div className="p-3 bg-base-200 rounded-lg">
                <span className="font-medium">
                  {monitoringMethods.find(m => m.value === monitoringConfig?.method)?.label}
                </span>
                <p className="text-sm text-gray-500 mt-1">
                  {monitoringMethods.find(m => m.value === monitoringConfig?.method)?.description}
                </p>
              </div>
            </div>

            <div>
              <label className="label">
                <span className="label-text">Host/URL</span>
              </label>
              <div className="p-3 bg-base-200 rounded-lg">
                <span>{monitoringConfig?.host}</span>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label className="label">
                  <span className="label-text">Max Failures</span>
                </label>
                <div className="p-3 bg-base-200 rounded-lg">
                  <span>{monitoringConfig?.max_failures}</span>
                </div>
              </div>

              <div>
                <label className="label">
                  <span className="label-text">Checking Interval</span>
                </label>
                <div className="p-3 bg-base-200 rounded-lg">
                  <span>{monitoringConfig?.checking_interval} seconds</span>
                </div>
              </div>

              <div>
                <label className="label">
                  <span className="label-text">Airplane Mode Delay</span>
                </label>
                <div className="p-3 bg-base-200 rounded-lg">
                  <span>{monitoringConfig?.airplane_mode_delay} seconds</span>
                </div>
              </div>
            </div>
            {monitoringConfig.is_active ? (
              <button
                onClick={onStop}
                disabled={isStopping}
                className="btn btn-error btn-sm w-full md:w-auto"
              >
                {isStopping ? (
                  <>
                    <span className="loading loading-spinner loading-xs"></span>
                    Stopping...
                  </>
                ) : (
                  <>
                    <FaCircleStop className="w-5 h-5 mr-2" />
                    Stop Monitoring
                  </>
                )}
              </button>
            ) : (
              <button
                onClick={onStart}
                disabled={isStarting || !allowed}
                className="btn btn-success btn-sm w-full md:w-auto"
              >
                {isStarting ? (
                  <>
                    <span className="loading loading-spinner loading-xs"></span>
                    Starting...
                  </>
                ) : (
                  <>
                    <IoIosPlayCircle className="w-5 h-5 mr-2" />
                    Start Monitoring
                  </>
                )}
              </button>
            )}
          </div>
        ) :
          <div className="text-center py-8">
            <div className="mb-4">
              <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <h3 className="text-lg font-medium text-base-900 mb-2">No Monitoring Configuration</h3>
            <p className="text-sm text-gray-500 mb-4">
              {
                allowed ? (
                  `No monitoring configuration has been set up for this device yet.
                  Create a configuration to start monitoring the device's connectivity.`
                ) : (
                  `You can't create monitoring config. because airplane mode feature not available. please check device feature availability.`
                )
              }
            </p>
            <button
              disabled={!allowed}
              onClick={onCreate}
              className="btn btn-primary"
            >
              Create Configuration
            </button>
          </div>
        }
      </div>
    </div>
  )
}

export default MonitoringConfigCard