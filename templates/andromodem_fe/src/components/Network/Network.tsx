import type { FC } from "react"
import { LuRefreshCw } from "react-icons/lu"

import useNetwork from "@/hooks/useNetwork"
import useTabManagement from "@/hooks/useTabManagement"
import { useDevicesStore } from "@/stores/devicesStore"
import { type Device,DeviceState } from "@/types/device"
import { isFeatureAvailable } from "@/utils/common"

import LoadingOverlay from "../Loading/LoadingOverlay"
import ApnModal from "../Modal/ApnModal"
import NetworkCardSkeleton from "../Skeleton/NetworkCardSkeleton"
import AccessPointName from "./AccessPointName"
import AirplaneMode from "./AirplaneMode"
import IpRoutes from "./IpRoutes"
import MobileData from "./MobileData"
import SimInfoDisplay from "./SimInfo"

const Network: FC<{ device: Device }> = ({ device }) => {
    const {
        networkInfo,
        isLoadingNetworkInfo,
        isValidatingNetworkInfo,
        mutateNetworkInfo,
        toggleMobileData,
        isTogglingMobileData,
        toggleAirplaneMode,
        isTogglingAirplaneMode
    } = useNetwork(device)
    const { activeTab, handleTabChange } = useTabManagement("1")

    const deviceFeatureAvailabilities  = useDevicesStore((state) => state.deviceFeatureAvailabilities)

    const isLoading = () => {
        return isLoadingNetworkInfo || isValidatingNetworkInfo || isTogglingMobileData || isTogglingAirplaneMode
    }
    if (networkInfo === null) {
        return (
            <NetworkCardSkeleton count={3} />
        )
    }

    return (
        <div>
            <button
                onClick={() => void mutateNetworkInfo()}
                disabled={
                    isValidatingNetworkInfo ||
                    isLoadingNetworkInfo ||
                    device == null ||
                    device.state === DeviceState.DISCONNECT
                }
                className="btn btn-xs md:btn-sm btn-active btn-primary mb-3"
            >
                <LuRefreshCw className={isValidatingNetworkInfo ? "animate-spin" : ""} />
                Refresh
            </button>

            {(isLoadingNetworkInfo || isValidatingNetworkInfo) && <LoadingOverlay />}

            {networkInfo.apn && <ApnModal apn={networkInfo.apn} modal_id="apn_modal" />}

            <div className="col-span-1 md:col-span-2 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5 md:gap-20 justify-items-center">
                <AirplaneMode
                    isAirplaneModeEnabled={networkInfo.airplane_mode}
                    isLoading={isLoading() || !isFeatureAvailable(deviceFeatureAvailabilities, "can_change_airplane_mode_status")}
                    device={device}
                    onToggle={toggleAirplaneMode}
                />

                <MobileData
                    sims={networkInfo.sims}
                    isLoading={isLoading()}
                    device={device}
                    disabled={networkInfo.airplane_mode}
                    onToggle={toggleMobileData}
                />

                <AccessPointName
                    apn={networkInfo.apn}
                    isLoading={isLoading()}
                    device={device}
                    modalId="apn_modal"
                />
            </div>

            <div className="col-span-1 md:col-span-2 grid grid-cols-1 lg:grid-cols-2 gap-5 md:gap-10 justify-items-center mt-10">
                <IpRoutes ipRoutes={networkInfo.ip_routes} />

                <SimInfoDisplay
                    sims={networkInfo.sims}
                    activeTab={activeTab}
                    onTabChange={handleTabChange}
                />
            </div>
        </div>
    )
}

export default Network