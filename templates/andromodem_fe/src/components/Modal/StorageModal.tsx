import type { DeviceStorage } from "@/types/device.ts"
import { convertStorageUnit } from "@/utils/converter"
type props = {
    storage: DeviceStorage;
    modal_id: string;
};
const ModalStorage = ({storage, modal_id}: props) => {
    return <dialog id={modal_id} className="modal">
        <div className="modal-box">
            <h3 className="font-bold text-lg">Storage</h3>
            <div className="my-5">
                <h5>Internal Storage</h5>
                <div className="flex justify-between">
                    <span>{convertStorageUnit((storage.system_used) + (storage.data_used), "KB", "GB")} GB Used</span>
                    <span>{convertStorageUnit((storage.system_total) + (storage.data_total), "KB", "GB")} GB Total</span>
                </div>
                <progress className="progress w-full"
                          value={(storage.system_used) + (storage.data_used)}
                          max={(storage.system_total) + (storage.data_total)}></progress>
            </div>
            <hr/>
            <div className="my-1">
                <h5>System</h5>
                <div className="flex justify-between">
                    <span>{convertStorageUnit(storage.system_used, "KB", "GB")} GB Used</span>
                    <span>{convertStorageUnit(storage.system_total, "KB", "GB")} GB Total</span>
                </div>
                <progress className="progress w-full" value={storage.system_used}
                          max={storage.system_total}></progress>
            </div>
            <div className="my-1">
                <h5>Data</h5>
                <div className="flex justify-between">
                    <span>{convertStorageUnit(storage.data_used, "KB", "GB")} GB Used</span>
                    <span>{convertStorageUnit(storage.data_total, "KB", "GB")} GB Total</span>
                </div>
                <progress className="progress w-full" value={storage.data_used}
                          max={storage.data_total}></progress>
            </div>
            <div className="modal-action">
                <form method="dialog">
                    <button className="btn btn-outline">Close</button>
                </form>
            </div>
        </div>
    </dialog>
}
export default ModalStorage