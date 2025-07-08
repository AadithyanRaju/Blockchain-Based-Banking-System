import tkinter as tk
from tkinter import messagebox
import json
import asyncio
from hfc.fabric import Client

# Hyperledger Fabric Client
loop = asyncio.get_event_loop()
c_hlf = Client(net_profile="path_to_connection_profile.json")

# Function to query balance
async def query_balance(username):
    try:
        user = c_hlf.get_user('org1.example.com', username)
        response = await c_hlf.chaincode_query(
            requestor=user,
            channel_name='mychannel',
            peer_names=['peer0.org1.example.com'],
            args=[username],
            cc_name='banking',
            fcn='queryBalance'
        )
        return response
    except Exception as e:
        return str(e)

# Function to transfer funds
async def transfer_funds(from_user, to_user, amount):
    try:
        user = c_hlf.get_user('org1.example.com', from_user)
        response = await c_hlf.chaincode_invoke(
            requestor=user,
            channel_name='mychannel',
            peer_names=['peer0.org1.example.com'],
            args=[from_user, to_user, amount],
            cc_name='banking',
            fcn='transferFunds'
        )
        return response
    except Exception as e:
        return str(e)

# Function to deposit funds
async def deposit_funds(username, amount):
    try:
        user = c_hlf.get_user('org1.example.com', username)
        response = await c_hlf.chaincode_invoke(
            requestor=user,
            channel_name='mychannel',
            peer_names=['peer0.org1.example.com'],
            args=[username, amount],
            cc_name='banking',
            fcn='depositFunds'
        )
        return response
    except Exception as e:
        return str(e)

# Function to withdraw funds
async def withdraw_funds(username, amount):
    try:
        user = c_hlf.get_user('org1.example.com', username)
        response = await c_hlf.chaincode_invoke(
            requestor=user,
            channel_name='mychannel',
            peer_names=['peer0.org1.example.com'],
            args=[username, amount],
            cc_name='banking',
            fcn='withdrawFunds'
        )
        return response
    except Exception as e:
        return str(e)

# GUI Application
class BlockchainApp:
    def __init__(self, root):
        self.root = root
        self.root.title("Blockchain Banking System")

        # Labels
        self.lbl_username = tk.Label(root, text="Username:")
        self.lbl_username.grid(row=0, column=0)
        self.lbl_amount = tk.Label(root, text="Amount:")
        self.lbl_amount.grid(row=1, column=0)
        self.lbl_recipient = tk.Label(root, text="Recipient:")
        self.lbl_recipient.grid(row=2, column=0)

        # Entries
        self.entry_username = tk.Entry(root)
        self.entry_username.grid(row=0, column=1)
        self.entry_amount = tk.Entry(root)
        self.entry_amount.grid(row=1, column=1)
        self.entry_recipient = tk.Entry(root)
        self.entry_recipient.grid(row=2, column=1)

        # Buttons
        self.btn_query = tk.Button(root, text="Query Balance", command=self.query_balance)
        self.btn_query.grid(row=3, column=0)
        self.btn_deposit = tk.Button(root, text="Deposit", command=self.deposit_funds)
        self.btn_deposit.grid(row=3, column=1)
        self.btn_withdraw = tk.Button(root, text="Withdraw", command=self.withdraw_funds)
        self.btn_withdraw.grid(row=4, column=0)
        self.btn_transfer = tk.Button(root, text="Transfer", command=self.transfer_funds)
        self.btn_transfer.grid(row=4, column=1)

    def query_balance(self):
        username = self.entry_username.get()
        result = loop.run_until_complete(query_balance(username))
        messagebox.showinfo("Query Balance", f"Balance for {username}: {result}")

    def deposit_funds(self):
        username = self.entry_username.get()
        amount = self.entry_amount.get()
        result = loop.run_until_complete(deposit_funds(username, amount))
        messagebox.showinfo("Deposit Funds", f"Deposit result: {result}")

    def withdraw_funds(self):
        username = self.entry_username.get()
        amount = self.entry_amount.get()
        result = loop.run_until_complete(withdraw_funds(username, amount))
        messagebox.showinfo("Withdraw Funds", f"Withdrawal result: {result}")

    def transfer_funds(self):
        from_user = self.entry_username.get()
        to_user = self.entry_recipient.get()
        amount = self.entry_amount.get()
        result = loop.run_until_complete(transfer_funds(from_user, to_user, amount))
        messagebox.showinfo("Transfer Funds", f"Transfer result: {result}")

# Run the GUI application
if __name__ == "__main__":
    root = tk.Tk()
    app = BlockchainApp(root)
    root.mainloop()
