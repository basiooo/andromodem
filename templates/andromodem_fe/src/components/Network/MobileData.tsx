import type { FC } from "react"
import { TbMobiledata, TbMobiledataOff } from "react-icons/tb"

import useMobileDataStatus from "@/hooks/useMobileDataStatus"
import type { Device } from "@/types/device"
import { DeviceState } from "@/types/device"
import type { SimInfo } from "@/types/network"

interface MobileDataProps {
  sims: SimInfo[]
  isLoading: boolean
  disabled: boolean
  device: Device
  onToggle?: () => void
}

const MobileData: FC<MobileDataProps> = ({
  sims,
  isLoading,
  disabled,
  device,
  onToggle
}) => {
  const { hasMobileDataEnabled } = useMobileDataStatus()
  const isMobileDataEnabled = hasMobileDataEnabled(sims)

  return (
    <div className="card bg-base-100 w-full shadow-sm">
      <div className="card-body">
        <h2 className="card-title">Mobile Data</h2>
        <div className="mx-auto flex flex-col text-center items-center">
          <div className="bg-base-300 rounded-full p-2">
            <div className="tooltip" data-tip={`Mobile data is ${isMobileDataEnabled ? "enabled" : "disabled"}`}>
              {
                isMobileDataEnabled ?
                  <TbMobiledata size="3rem" /> :
                  <TbMobiledataOff size="3rem" />
              }
            </div>
          </div>
          <span className="text-sm text-neutral-content/50">
          Requires further verification, may not work
          </span>
          <button
            className={`mt-5 btn btn-sm ${isMobileDataEnabled ? "btn-error" : "btn-success"}`}
            disabled={isLoading || device == null || device.state === DeviceState.DISCONNECT || disabled}
            onClick={onToggle}
          >
            {isMobileDataEnabled ? "Disable" : "Enable"}
          </button>
        </div>
      </div>
    </div>
  )
}

export default MobileData