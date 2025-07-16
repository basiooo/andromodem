import type { FC } from "react"
import {FaDonate, FaFacebook, FaGithub, FaHeart, FaInstagram} from "react-icons/fa"
import {FaEarthAsia} from "react-icons/fa6"

import {config} from "@/config"


const Footer: FC = () => {
    return (
        <footer className="footer sm:footer-horizontal bg-base-300 text-base p-10">
            <div className="container flex justify-between footer mx-auto flex-col md:flex-row">
                <aside className="mx-auto md:mx-0">
                    <b className="mx-auto md:mx-0">AndroModem version {config.VERSION}</b>
                    <p className="mx-auto md:mx-0">
                        Made with
                        <FaHeart
                            className="inline text-xl mx-2 animate-pulse animate-infinite animate-duration-500 animate-ease-in text-pink-600"/>
                        in Nganjuk
                    </p>
                    <p className="text-md mx-auto md:mx-0">2025 @ <a className="link" target="_blank" rel="noreferrer"
                                                                     href="https://github.com/basiooo">Bagas
                        Julianto</a></p>
                </aside>
                <nav
                    className="grid-flow-col gap-4 md:place-self-center md:justify-self-end mx-auto md:mx-0 mt-5 md:mt-0">
                    <a className="link" target="_blank" rel="noreferrer" href="https://github.com/basiooo/andromodem">
                        <FaGithub className="inline mx-1 text-2xl"/></a>
                    <a className="link" target="_blank" rel="noreferrer" href="https://saweria.co/basiooo">
                        <FaDonate className="inline mx-1 text-2xl"/></a>
                    <a className="link" target="_blank" rel="noreferrer" href="https://bagasjulianto.my.id">
                        <FaEarthAsia className="inline mx-1 text-2xl"/></a>
                    <a className="link" target="_blank" rel="noreferrer" href="https://www.facebook.com/bagas.jul.33">
                        <FaFacebook className="inline mx-1 text-2xl"/></a>
                    <a className="link" target="_blank" rel="noreferrer" href="https://www.instagram.com/_basiooo">
                        <FaInstagram className="inline mx-1 text-2xl"/></a>
                </nav>
            </div>
        </footer>
    )
}
export default Footer