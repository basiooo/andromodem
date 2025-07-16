import type { FC } from "react"
import { MdAirplanemodeActive, MdAirplanemodeInactive } from "react-icons/md"

import type { Device } from "@/types/device"
import { DeviceState } from "@/types/device"

interface AirplaneModeProps {
  isAirplaneModeEnabled: boolean
  isLoading: boolean
  device: Device
  onToggle?: () => void
}

const AirplaneMode: FC<AirplaneModeProps> = ({
  isAirplaneModeEnabled,
  isLoading,
  device,
  onToggle
}) => {
  return (
    <div className="card bg-base-100 w-full shadow-sm">
      <div className="card-body">
        <h2 className="card-title">Airplane Mode</h2>
        <div className="mx-auto flex flex-col text-center items-center">
          <div className="bg-base-300 rounded-full p-2">
            <div className="tooltip" data-tip={`Airplane mode is ${isAirplaneModeEnabled ? "enabled" : "disabled"}`}>
              {
                isAirplaneModeEnabled ?
                  <MdAirplanemodeActive size="3rem" /> :
                  <MdAirplanemodeInactive size="3rem" />
              }
            </div>
          </div>
          <button
            className={`mt-5 btn btn-sm ${isAirplaneModeEnabled ? "btn-success" : "btn-error"}`}
            disabled={isLoading || device == null || device.state === DeviceState.DISCONNECT}
            onClick={onToggle}
          >
            {isAirplaneModeEnabled ? "Disable" : "Enable"}
          </button>
        </div>
      </div>
    </div>
  )
}

export default AirplaneMode