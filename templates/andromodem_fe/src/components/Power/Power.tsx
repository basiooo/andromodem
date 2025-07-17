import { FaAndroid, FaPowerOff } from "react-icons/fa"
import { IoMdSettings } from "react-icons/io"
import { MdOutlineRestartAlt } from "react-icons/md"

import PowerActionConfirmModal from "@/components/Modal/PowerActionConfirmModal"
import usePowerActions from "@/hooks/usePowerActions"
import { type Device } from "@/types/device"


type props = {
    device: Device;
};

const Power = ({ device }: props) => {
    const { isDisabled, powerActionHandler, powerActionHandlerConfirm } = usePowerActions(device)
    return (
        <>
            <PowerActionConfirmModal onConfirm={powerActionHandlerConfirm} modal_id="power_modal" />
            <div className="card card-compact w-full bg-base-200 shadow-xl">
                <div className="card-body">
                    <div className="grid grid-cols-2 justify-items-center gap-10">
                        <div className="text-center w-fit">
                            <button disabled={isDisabled} className="btn w-16 h-16 md:w-24 md:h-24 btn-circle bg-red-500 hover:bg-red-700" onClick={() => powerActionHandler("power_off")}>
                                <FaPowerOff className="text-2xl md:text-4xl" />
                            </button>
                            <p className="text-base md:text-xl">Power Off</p>
                        </div>
                        <div className="text-center w-fit">
                            <button disabled={isDisabled} className="btn w-16 h-16 md:w-24 md:h-24 btn-circle bg-red-500 hover:bg-red-700" onClick={() => powerActionHandler("reboot")}>
                                <MdOutlineRestartAlt className="text-2xl md:text-4xl" />
                            </button>
                            <p className="text-base md:text-xl">Reboot</p>
                        </div>
                        <div className="text-center w-fit">
                            <button disabled={isDisabled} className="btn w-16 h-16 md:w-24 md:h-24 btn-circle bg-red-500 hover:bg-red-700" onClick={() => powerActionHandler("reboot_recovery")}>
                                <IoMdSettings className="text-2xl md:text-4xl" />
                            </button>
                            <p className="text-base md:text-xl">Reboot Recovery</p>
                        </div>
                        <div className="text-center w-fit">
                            <button disabled={isDisabled} className="btn w-16 h-16 md:w-24 md:h-24 btn-circle bg-red-500 hover:bg-red-700" onClick={() => powerActionHandler("reboot_bootloader")}>
                                <FaAndroid className="text-2xl md:text-4xl" />
                            </button>
                            <p className="text-base md:text-xl">Reboot Bootloader</p>
                        </div>
                    </div>
                </div>
            </div>
        </>
    )
}

export default Power