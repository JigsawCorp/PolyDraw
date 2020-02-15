﻿using System;
using System.ComponentModel;
using System.Runtime.CompilerServices;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Input;
using ClientLourd.Annotations;
using ClientLourd.Models.Bindable;
using ClientLourd.Utilities.Commands;
using MaterialDesignThemes.Wpf;

namespace ClientLourd.Views.Dialogs
{
    public partial class RegisterDialog : UserControl, INotifyPropertyChanged
    {
        public PrivateProfileInfo PrivateProfileInfo { get; set; }
        public RegisterDialog(PrivateProfileInfo infos)
        {
            PrivateProfileInfo = infos;
            Avatar = "";
            InitializeComponent();
        }

        public bool IsPasswordInvalid
        {
            get { return CheckInvalidPassword(); }
        }

        public string Avatar { get; set; }

        private bool CheckInvalidPassword()
        {
            if (PasswordField1.Password != PasswordField2.Password)
                return true;
            if (String.IsNullOrWhiteSpace(PasswordField1.Password) || PasswordField1.Password.Length < 8)
                return true;
            if (String.IsNullOrWhiteSpace(PasswordField2.Password) || PasswordField2.Password.Length < 8)
                return true;
            return false;
        }

        public event PropertyChangedEventHandler PropertyChanged;

        [NotifyPropertyChangedInvocator]
        protected virtual void OnPropertyChanged([CallerMemberName] string propertyName = null)
        {
            PropertyChanged?.Invoke(this, new PropertyChangedEventArgs(propertyName));
        }

        private void OnPasswordChanged(object sender, RoutedEventArgs e)
        {
            OnPropertyChanged(nameof(IsPasswordInvalid));
        }

        RelayCommand<Channel> _changeAvatarCommand;

        public ICommand ChangeAvatarCommand
        {
            get
            {
                return _changeAvatarCommand ??
                       (_changeAvatarCommand = new RelayCommand<Channel>(channel => ChangeAvatar()));
            }
        }

        private async void ChangeAvatar()
        {
            await DialogHost.Show(new AvatarSelectionDialog(), "RegisterDialogHost");
        }

    }
}