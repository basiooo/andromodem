import { type FC } from "react"
import { FaFingerprint, FaLinux, FaShieldAlt } from "react-icons/fa"
import { FaBoxOpen } from "react-icons/fa6"
import { GiAndroidMask, GiGearHammer } from "react-icons/gi"
import {
    IoHardwareChipOutline,
    IoPhonePortraitOutline,
    IoTimeSharp
} from "react-icons/io5"
import { LiaAndroid } from "react-icons/lia"
import { LuRefreshCw } from "react-icons/lu"
import { MdOutlineFactory } from "react-icons/md"
import { PiIdentificationBadgeBold } from "react-icons/pi"
import {TbBatteryVerticalOff} from "react-icons/tb"

import useDeviceInfo from "@/hooks/useDeviceInfo"
import { type Device, DeviceState } from "@/types/device"
import { convertStorageUnit, getPercentage, secondsForHuman } from "@/utils/converter"

import LoadingOverlay from "@/components/Loading/LoadingOverlay"
import ModalBattery from "@/components/Modal/BatteryModal"
import ModalMemory from "@/components/Modal/MemoryModal"
import ModalStorage from "@/components/Modal/StorageModal"
import DeviceInfoCardSkeleton from "@/components/Skeleton/DeviceInfoCardSkeleton"
import DeviceInfoItemSkeleton from "@/components/Skeleton/DeviceInfoItemSkeleton"
import DeviceInfoCard from "@/components/DeviceInfo/DeviceInfoCard"
import DeviceInfoItem from "@/components/DeviceInfo/DeviceInfoItem"

const DeviceInfo: FC<{ device: Device }> = ({ device }) => {
    const { deviceInfo, upTimeSecond, isLoading, isValidating, mutate } = useDeviceInfo(device)

    if (deviceInfo == null) {
        return (
            <>
            <DeviceInfoCardSkeleton />
            <DeviceInfoItemSkeleton/>
            </>
        )
    }
    return (
        <div>
            <button
                onClick={() => void mutate()}
                disabled={
                    isValidating ||
                    isLoading ||
                    device == null ||
                    device.state === DeviceState.DISCONNECT
                }
                className="btn btn-xs md:btn-sm btn-active btn-primary mb-3"
            >
                <LuRefreshCw className={isValidating ? "animate-spin" : ""} />
                Refresh
            </button>
            {isLoading || isValidating ? (
                <LoadingOverlay />
            ) : (
                <></>
            )}
            <ModalMemory modal_id="modal_memory" memory={deviceInfo?.memory}/>
            <ModalStorage modal_id="modal_storage" storage={deviceInfo?.storage}/>
            <ModalBattery modal_id="modal_battery" battery={deviceInfo.battery}/>
            <div
                className="col-span-1 md:col-span-2 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5 md:gap-10  justify-items-center">
                <DeviceInfoCard title="Storage" reverse_radial_color={true} modal_id="modal_storage"
                    value={getPercentage(deviceInfo.storage.data_total + deviceInfo.storage.system_total, deviceInfo.storage.data_used + deviceInfo.storage.system_used)}>
                    <p>{convertStorageUnit(deviceInfo.storage.data_used + deviceInfo.storage.system_used, "KB", "GB")} GB
                        Used</p>
                    <p>{convertStorageUnit(deviceInfo.storage.data_total + deviceInfo.storage.system_total, "KB", "GB")} GB
                        Total</p>
                </DeviceInfoCard>

                <DeviceInfoCard title="Memory" reverse_radial_color={true} modal_id="modal_memory"
                    value={getPercentage(
                        deviceInfo.memory.mem_total + deviceInfo.memory.swap_total,
                        deviceInfo.memory.mem_used + deviceInfo.memory.swap_used)}>
                    <p>{convertStorageUnit(deviceInfo.memory.mem_used + deviceInfo.memory.swap_used, "KB", "GB")} GB Used</p>
                    <p>{convertStorageUnit(deviceInfo.memory.mem_total + deviceInfo.memory.swap_total, "KB", "GB")} GB Total</p>
                </DeviceInfoCard>
                {
                    deviceInfo.battery.present ?
                    <DeviceInfoCard title="Battery" modal_id="modal_battery"
                                    value={getPercentage(100, deviceInfo.battery.level)}>
                        <p>{deviceInfo.battery.level} %</p>
                        <p>{deviceInfo.battery.temperature} °C</p>
                        <p>{deviceInfo.battery.status}</p>
                    </DeviceInfoCard> :
                    <div className="card bg-base-200 w-full shadow-sm">
                        <div className="card-body">
                            <h2 className="card-title">Battery</h2>
                            <div className="flex justify-center items-center flex-wrap flex-col">
                                <TbBatteryVerticalOff size="5rem" color="red"/>
                                <div className="text-xs sm:xs md:text-sm">
                                    <b>Device does not have a battery</b>
                                </div>
                            </div>
                        </div>
                    </div>
                }
                <div className="card bg-base-100 w-full shadow-sm">
                    <div className="card-body">
                        <h2 className="card-title">Root</h2>
                        <div className="mx-auto flex flex-col text-center items-center">
                            {deviceInfo.root.is_rooted ? (
                                <div className="flex items-center justify-around">
                                    <div className="bg-base-300 rounded-full p-2">
                                        <GiAndroidMask size="5rem" color="#05df73"/>
                                    </div>
                                    <div className="md:text-xs">
                                        {deviceInfo.root.name}:{deviceInfo.root.version}
                                    </div>
                                </div>
                            ) : (
                                <>
                                    <div className="bg-base-300 rounded-full p-2">
                                        <GiAndroidMask size={80} color="#FF3D00"/>
                                    </div>
                                    Device Not Rooted
                                </>
                            )}
                        </div>
                    </div>
                </div>
            </div>
            <ul className="list bg-base-200 rounded-box shadow-md my-5">
                <li className="p-4 pb-2 text-xl tracking-wide">Device</li>
                <DeviceInfoItem icon={<MdOutlineFactory size={30} />} label="Brand" value={deviceInfo.prop.brand} />
                <DeviceInfoItem icon={<IoPhonePortraitOutline size={25} />} label="Model" value={deviceInfo.prop.model} />
                <DeviceInfoItem icon={<PiIdentificationBadgeBold size={25} />} label="Name" value={deviceInfo.prop.name} />
                <DeviceInfoItem icon={<LiaAndroid size={25} />} label="Android Version" value={deviceInfo.prop.android_version} />
                <DeviceInfoItem icon={<FaLinux size={25} />} label="Kernel Version" value={deviceInfo.other.kernel_version} />
                <DeviceInfoItem icon={<FaFingerprint size={25} />} label="Fingerprint" value={deviceInfo.prop.fingerprint} />
                <DeviceInfoItem icon={<FaShieldAlt size={25} />} label="Security Patch" value={deviceInfo.prop.security_patch} />
                <DeviceInfoItem icon={<IoHardwareChipOutline size={25} />} label="Processor" value={`${deviceInfo.prop.processor} (${deviceInfo.prop.abi})`} />
                <DeviceInfoItem icon={<GiGearHammer size={25} />} label="SDK/API Level" value={deviceInfo.prop.sdk} />
                <DeviceInfoItem icon={<IoTimeSharp size={25} />} label="Uptime" value={secondsForHuman(upTimeSecond)} />
                <DeviceInfoItem icon={<FaBoxOpen size={25} />} label="Busybox Installed" value={deviceInfo.other.busybox_installed ? "Yes" : "No"} />
            </ul>
        </div>
    )
}
export default DeviceInfo