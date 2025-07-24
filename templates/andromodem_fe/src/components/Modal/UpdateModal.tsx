import Markdown from 'react-markdown'
import remarkBreaks from 'remark-breaks'
import remarkGfm from 'remark-gfm'

import type { UpdateInfo } from "@/types/update"

type Props = {
    updateInfo: UpdateInfo | null;
    modal_id: string;
};

const UpdateModal = ({ updateInfo, modal_id }: Props) => {
    if (!updateInfo) return null

    const formatReleaseDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString("en", {
            year: "numeric",
            month: "long",
            day: "numeric",
            hour: "2-digit",
            minute: "2-digit"
        })
    }

    return (
        <dialog id={modal_id} className="modal">
            <div className="modal-box max-w-4xl">
                <h3 className="font-bold text-lg">Update Information</h3>

                <div className="grid grid-cols-1 mt-5 md:grid-cols-2 gap-4 mb-6">
                    <div className="bg-base-300 p-4 rounded-lg">
                        <h4 className="font-semibold text-md md:text-lg mb-2">Current Version</h4>
                        <p className="text-sm md:text-lg font-mono">{updateInfo.currentVersion}</p>
                    </div>
                    <div className="bg-base-300 p-4 rounded-lg">
                        <h4 className="font-semibold text-md md:text-lg mb-2">Latest Version</h4>
                        <p className="text-sm md:text-lg font-mono">{updateInfo.latestVersion}</p>
                    </div>
                </div>

                {updateInfo.releaseInfo && (
                    <>
                        <div className="mb-6">
                            <div className="flex items-center justify-between mb-3">
                                <h4 className="font-semibold text-md md:text-lg">{updateInfo.releaseInfo.name}</h4>
                                <span className="text-sm text-gray-500">
                                    {formatReleaseDate(updateInfo.releaseInfo.published_at)}
                                </span>
                            </div>

                            <div className="bg-base-300 p-4 rounded-lg max-h-96 overflow-y-auto">
                                <h5 className="font-semibold mb-3">Changelog:</h5>
                                <div className="prose prose-sm max-w-none dark:prose-invert">
                                    <Markdown 
                                        remarkPlugins={[remarkGfm, remarkBreaks]}
                                        components={{
                                            h1: ({children}) => <h1 className="text-lg font-bold mb-2">{children}</h1>,
                                            h2: ({children}) => <h2 className="text-md font-semibold mb-2">{children}</h2>,
                                            h3: ({children}) => <h3 className="text-sm font-semibold mb-1">{children}</h3>,
                                            p: ({children}) => <p className="mb-2 leading-relaxed">{children}</p>,
                                            ul: ({children}) => <ul className="list-disc list-inside mb-2 space-y-1">{children}</ul>,
                                            li: ({children}) => <li className="ml-2">{children}</li>,
                                            a: ({href, children}) => (
                                                <a 
                                                    href={href} 
                                                    target="_blank" 
                                                    rel="noopener noreferrer" 
                                                    className="text-blue-500 hover:text-blue-700 underline"
                                                >
                                                    {children}
                                                </a>
                                            ),
                                            code: ({children}) => (
                                                <code className="bg-gray-200 dark:bg-gray-700 px-1 py-0.5 rounded text-xs font-mono">
                                                    {children}
                                                </code>
                                            )
                                        }}
                                    >
                                        {updateInfo.releaseInfo.body || ""}
                                    </Markdown>
                                </div>
                            </div>
                        </div>
                        <div className="mb-6">
                            <div className="flex items-center justify-between mb-3">
                                <h4 className="font-semibold text-lg">How to Update</h4>
                            </div>
                            <a className="link" target="_blank" rel="noreferrer" href="https://github.com/basiooo/andromodem?tab=readme-ov-file#option-3-openwrt-installation">See update instructions.</a>
                        </div>
                    </>
                )}

                <div className="modal-action">
                    <form method="dialog">
                        <button className="btn btn-outline" id={`close_${modal_id}`}>Close</button>
                    </form>
                </div>
            </div>
        </dialog>
    )
}

export default UpdateModal