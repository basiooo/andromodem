import type { Apn } from "@/types/network"

type props = {
    apn: Apn;
    modal_id: string;
};
const ApnModal = ({ apn, modal_id }: props) => {
    return <dialog id={modal_id} className="modal">
        <div className="modal-box">
            <h3 className="font-bold text-lg">Access Point Name</h3>
            <p className="text-sm text-gray-500 mb-4">
                Access Point Name is the name of the network that the device is connected to.
            </p>
            <ul className="list text-xs">
                {
                    Object.entries(apn).map(([key, value]) => {
                        return (
                            <li key={key} className="list-row">
                                <div>
                                    <div className="font-semibold">{key}</div>
                                    <div className="text-xs font-semibold opacity-60">
                                        {value ? value : "-"}
                                    </div>
                                </div>
                            </li>
                        )
                    })
                }
            </ul>
            <div className="modal-action">
                <form method="dialog">
                    <button className="btn btn-outline">Close</button>
                </form>
            </div>
        </div>
    </dialog>
}
export default ApnModal
