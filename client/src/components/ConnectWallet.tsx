import { useState, useEffect } from "react";
import { ethers } from "ethers";
import "./ConnectWallet.css";
import { GotchiSelectModal } from "./GotchiSelectModal";
import gotchiIcon from "../assets/gotchi-icon.png"; // Placeholder
import walletIcon from "../assets/wallet-icon.png"; // Placeholder
import { FaWallet } from "react-icons/fa";

export const ConnectWallet = () => {
    const [provider, setProvider] = useState<ethers.BrowserProvider | null>(null);
    const [account, setAccount] = useState<string | null>(null);
    const [network, setNetwork] = useState<string | null>(null);
    const [showModal, setShowModal] = useState(false);

    useEffect(() => {
        checkWalletConnection();
    }, []);

    const checkWalletConnection = async () => {
        if (window.ethereum) {
            const web3Provider = new ethers.BrowserProvider(window.ethereum);
            const accounts = await web3Provider.listAccounts();
            if (accounts.length > 0) {
                const address = await accounts[0].getAddress(); // Extract address
                setProvider(web3Provider);
                setAccount(address);
                const net = await web3Provider.getNetwork();
                setNetwork(net.name === "matic" ? "Polygon" : net.name);
            }
        }
    };

    const connectWallet = async () => {
        if (!window.ethereum) {
            alert("Please install a Web3 wallet like MetaMask!");
            return;
        }

        const web3Provider = new ethers.BrowserProvider(window.ethereum);
        try {
            await window.ethereum.request({ method: "eth_requestAccounts" });
            const accounts = await web3Provider.listAccounts();
            const address = await accounts[0].getAddress(); // Extract address
            setProvider(web3Provider);
            setAccount(address);

            const net = await web3Provider.getNetwork();
            if (net.chainId !== 137n) { // Polygon chain ID
                await window.ethereum.request({
                    method: "wallet_switchEthereumChain",
                    params: [{ chainId: "0x89" }],
                });
            }
            setNetwork("Polygon");
        } catch (error) {
            console.error("Wallet connection failed:", error);
        }
    };

    const formatAddress = (addr: string) =>
        `${addr.slice(0, 4)}...${addr.slice(-4)}`;

    return (
        <div className="connect-wallet-container">
            {!account ? (
                <button className="connect-button" onClick={connectWallet}>
                    Connect Wallet
                </button>
            ) : (
                <div className="wallet-bar">
                    <img src={gotchiIcon} alt="Gotchi" className="gotchi-icon" />
                    <span
                        className="find-gotchi"
                        onClick={() => setShowModal(true)}
                    >
                        Find your Gotchi
                    </span>
                    <span style={{width:"32px"}} />
                    <span className="network-name">{network}</span>
                    <FaWallet className="wallet-icon" />
                    <span className="wallet-address">
                        {formatAddress(account)}
                    </span>
                </div>
            )}
            {showModal && account && (
                <GotchiSelectModal
                    account={account}
                    onClose={() => setShowModal(false)}
                />
            )}
        </div>
    );
};