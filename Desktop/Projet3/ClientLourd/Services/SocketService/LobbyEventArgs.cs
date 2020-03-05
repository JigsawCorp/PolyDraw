﻿using ClientLourd.Models.Bindable;
using System;
using System.Collections.Generic;
using System.Collections.ObjectModel;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ClientLourd.Services.SocketService
{
    public class LobbyEventArgs: EventArgs
    {
        private dynamic _data;
        public LobbyEventArgs(dynamic data)
        {
            _data = data;
            
        }

        // Lobby created
        public string ID { get => _data["ID"]; }

        public string Name { get => _data["Name"]; }

        public string OwnerName { get => _data["OwnerName"]; }

        public string OwnerID { get => _data["OwnerID"]; }

        public int PlayersMax { get => (int)_data["PlayersMax"]; }

        public int Mode { get => (int)_data["Mode"]; }

        public int Difficulty { get => (int)_data["Difficulty"]; }


        public ObservableCollection<Player> Players { 
            get 
            {
                ObservableCollection<Player> players = new ObservableCollection<Player>();


                for (int i = 0; i < ((dynamic[])_data["Players"]).Length; i++)
                {
                    players.Add(new Player(_data["Players"][i]["IsCPU"], _data["Players"][i]["ID"], _data["Players"][i]["Username"]));
                }
                return players;
            }
        }

        // Join Lobby Response
        public bool Response
        {
            get { return _data["Response"]; }
        }
        public string Error
        {
            get { return _data["Error"]; }
        }

        // UserJoinedLobby and QuitLobbyResponse

        public bool IsCPU
        {
            get { return _data["IsCPU"]; }
        }
        public string UserID
        {
            get { return _data["UserID"]; }
        }
        public string Username
        {
            get { return _data["Username"]; }
        }
        public string GroupID
        {
            get { return _data["GroupID"]; }
        }


    }
}
