import type { FC } from "react"

import { type MonitoringFormData,useMonitoringForm } from "@/hooks/useMonitoringForm"
import type { MonitoringConfig } from "@/types/monitoring"
import { type MonitoringMethodValue } from "@/types/monitoring"
import { monitoringMethods } from "@/utils/monitoringUtils"

interface MonitoringConfigModalProps {
  modalId: string;
  monitoringConfig: MonitoringConfig | null;
  onSubmit: (data: MonitoringFormData) => Promise<void>;
  onCancel: () => void;
}

const MonitoringConfigModal: FC<MonitoringConfigModalProps> = ({
  modalId,
  monitoringConfig,
  onSubmit,
  onCancel
}) => {
  const {
    register,
    handleSubmit,
    errors,
    isSubmitting,
    setCurrentMethod,
    handleCancel: handleFormCancel
  } = useMonitoringForm({
    monitoringConfig,
    onSubmit,
    onCancel
  })

  return (
    <dialog id={modalId} className="modal">
      <div className="modal-box max-w-2xl">
        <h3 className="font-bold text-lg mb-4">
          {monitoringConfig ? 'Edit Monitoring Configuration' : 'Create Monitoring Configuration'}
        </h3>
        
        <form onSubmit={handleSubmit}>
          <div className="form-control w-full">
            <label className="label">
              <span className="label-text font-semibold">Monitoring Method</span>
            </label>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3 mt-2">
              {monitoringMethods.map((method) => (
                <div key={method.value} className="form-control">
                  <label className="label cursor-pointer w-full border rounded-lg p-3 hover:bg-base-200 transition-colors">
                    <div className="flex-1">
                      <div className="flex items-center gap-3">
                        <input
                          type="radio"
                          className="radio radio-primary"
                          value={method.value}
                          {...register('method', {
                            onChange: (e) => {
                              const newMethod = e.target.value as MonitoringMethodValue
                              setCurrentMethod(newMethod)
                            }
                          })}
                        />
                        <div>
                          <span className="label-text sm:text-sm font-medium">{method.label}</span>
                          <p className="text-xs text-gray-500 mt-1 text-wrap">{method.description}</p>
                        </div>
                      </div>
                    </div>
                  </label>
                </div>
              ))}
            </div>
            {errors.method && (
              <label className="label">
                <span className="text-xs text-balance text-error">{errors.method.message}</span>
              </label>
            )}
          </div>

          <div className="form-control w-full mt-4">
            <label className="label block">
              <span className="label-text">Host/URL</span>
            </label>
            <label className="input validator w-full">
              <input
                type="text"
                placeholder="example.com"
                className={`w-full ${errors.host ? 'input-error' : ''}`}
                {...register('host')}
              />
            </label>
            <div className="label">
              <span className="text-xs text-balance text-gray-500">Enter domain or IP address without protocol (http/https)</span>
            </div>
            {errors.host && (
              <label className="label block">
                <span className="text-xs text-balance text-error">{errors.host.message}</span>
              </label>
            )}
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4">
            <div className="form-control w-full">
              <label className="label">
                <span className="label-text">Max Failures</span>
              </label>
              <input
                type="number"
                placeholder="3"
                className={`input input-bordered w-full ${errors.max_failures ? 'input-error' : ''}`}
                {...register('max_failures', { valueAsNumber: true })}
              />
              <div className="label">
                <span className="text-xs text-balance text-gray-500">Maximum number of failures before activating airplane mode (1-100)</span>
              </div>
              {errors.max_failures && (
                <label className="label">
                  <span className="text-xs text-balance text-error">{errors.max_failures.message}</span>
                </label>
              )}
            </div>

            <div className="form-control w-full">
              <label className="label">
                <span className="label-text">Checking Interval (seconds)</span>
              </label>
              <input
                type="number"
                placeholder="60"
                className={`input input-bordered w-full ${errors.checking_interval ? 'input-error' : ''}`}
                {...register('checking_interval', { valueAsNumber: true })}
              />
              <div className="label">
                <span className="text-xs text-balance text-gray-500">Time interval for connection checking (5-3600 seconds)</span>
              </div>
              {errors.checking_interval && (
                <label className="label">
                  <span className="text-xs text-balance text-error">{errors.checking_interval.message}</span>
                </label>
              )}
            </div>

            <div className="form-control w-full">
              <label className="label">
                <span className="label-text">Airplane Mode Delay (seconds)</span>
              </label>
              <input
                type="number"
                placeholder="1"
                min={1}
                className={`input input-bordered w-full ${errors.airplane_mode_delay ? 'input-error' : ''}`}
                {...register('airplane_mode_delay', { valueAsNumber: true })}
              />
              <div className="label">
                <span className="text-xs text-balance text-gray-500">Delay time before activating airplane mode after failure detection (1-3600 seconds)</span>
              </div>
              {errors.airplane_mode_delay && (
                <label className="label">
                  <span className="text-xs text-balance text-error">{errors.airplane_mode_delay.message}</span>
                </label>
              )}
            </div>
          </div>

          <div className="modal-action">
            <button
              type="button"
              onClick={handleFormCancel}
              className="btn btn-outline"
              disabled={isSubmitting}
            >
              Cancel
            </button>
            <button
              type="submit"
              className={`btn btn-primary ${isSubmitting ? 'loading' : ''}`}
              disabled={isSubmitting}
            >
              {isSubmitting ? 'Saving...' : (monitoringConfig ? 'Save Configuration' : 'Create Configuration')}
            </button>
          </div>
        </form>
      </div>
      <form method="dialog" className="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>
  )
}

export default MonitoringConfigModal