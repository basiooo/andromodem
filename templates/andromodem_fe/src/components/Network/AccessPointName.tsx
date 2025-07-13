import type { FC } from "react"
import { MdPermDataSetting } from "react-icons/md"

import type { Device } from "@/types/device"
import { DeviceState } from "@/types/device"
import type { Apn } from "@/types/network"
import { showModal } from "@/utils/common"

interface AccessPointNameProps {
  apn: Apn | null
  isLoading: boolean
  device: Device
  modalId: string
}

const AccessPointName: FC<AccessPointNameProps> = ({
  apn,
  isLoading,
  device,
  modalId
}) => {
  return (
    <div className="card bg-base-100 w-full shadow-sm">
      <div className="card-body">
        <h2 className="card-title">Access Point Name</h2>
        <div className="mx-auto flex flex-col text-center items-center">
          <div className="bg-base-300 rounded-full p-2">
            <MdPermDataSetting size="3rem" />
          </div>
          <button
            className={`mt-5 btn btn-sm btn-info`}
            onClick={() => showModal(modalId)}
            disabled={apn === null || isLoading || device == null || device.state === DeviceState.DISCONNECT}
          >
            Detail
          </button>
        </div>
      </div>
    </div>
  )
}

export default AccessPointName