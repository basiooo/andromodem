import type { DeviceMemory } from "@/types/device.ts"
import { convertStorageUnit } from "@/utils/converter"

type props = {
    memory: DeviceMemory;
    modal_id: string;
};
const ModalMemory = ({memory, modal_id}: props) => {
    return <dialog id={modal_id} className="modal">
        <div className="modal-box">
            <h3 className="font-bold text-lg">Memory</h3>
            <div className="my-5">
                <div className="flex justify-between">
                    <span>RAM</span>
                    <span>{convertStorageUnit(memory.mem_total, "KB", "GB")} GB</span>
                </div>
                <div className="flex justify-between">
                    <span>{convertStorageUnit(memory.mem_used, "KB", "GB")} GB Used</span>
                    <span>{convertStorageUnit(memory.mem_free, "KB", "GB")} GB Free</span>
                </div>
                <progress className="progress w-full" value={memory.mem_used} max={memory.mem_total}></progress>
            </div>
            {
                memory.swap_total ?
                    <div className="my-5">
                        <div className="flex justify-between">
                            <span>Extended RAM</span>
                            <span>{convertStorageUnit(memory.swap_total, "KB", "GB")} GB</span>
                        </div>
                        <div className="flex justify-between">
                            <span>{convertStorageUnit(memory.swap_used, "KB", "GB")} GB Used</span>
                            <span>{convertStorageUnit(memory.swap_free, "KB", "GB")} GB Free</span>
                        </div>
                        <progress className="progress w-full" value={memory.swap_used}
                                  max={memory.swap_total}></progress>
                    </div>
                    : <></>
            }
            <div className="modal-action">
                <form method="dialog">
                    <button className="btn btn-outline">Close</button>
                </form>
            </div>
        </div>
    </dialog>
}
export default ModalMemory