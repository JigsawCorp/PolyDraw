﻿using System;
using System.ComponentModel;
using System.Runtime.CompilerServices;
using System.Text.RegularExpressions;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Input;
using System.Windows.Media.Imaging;
using ClientLourd.Annotations;
using ClientLourd.Models.Bindable;
using ClientLourd.Utilities.Commands;
using MaterialDesignThemes.Wpf;

namespace ClientLourd.Views.Dialogs
{
    public partial class RegisterDialog : UserControl, INotifyPropertyChanged
    {
        public User User { get; set; }
        public RegisterDialog(User user)
        {
            User = user;
            InitializeComponent();
        }
        public bool IsPasswordInvalid
        {
            get { return CheckInvalidPassword(); }
        }

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
            var result = await DialogHost.Show(new AvatarSelectionDialog(), "RegisterDialogHost");
            User.Avatar = (BitmapImage) result;
        }

    }
}