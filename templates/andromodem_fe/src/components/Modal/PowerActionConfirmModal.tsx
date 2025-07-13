type props = {
    onConfirm: () => Promise<void>;
    modal_id: string;
};
const PowerActionConfirmModal = ({onConfirm, modal_id}: props) => {
    return <dialog id={modal_id} className="modal backdrop-blur-sm">
        <div className="modal-box text-center">
            <h3 className="font-bold text-3xl">Confirm</h3>
            <p className="text-xl py-4">Are you sure?</p>
            <div className="modal-action">
                <form method="dialog">
                    <button className="btn btn-outline" id={`close_${modal_id}`}>Close</button>
                </form>
                <button className="btn btn-primary" onClick={onConfirm}>Confirm</button>
            </div>
        </div>
    </dialog>
}
export default PowerActionConfirmModal