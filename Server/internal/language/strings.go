package language

var strings = []byte(`
en:
  error:
    channelNotFound: "The specified channel cannot be found."
    channelInvalidStart: "Invalid parameters, start must be the lowest parameter."
    channelInvalidUrl: "Invalid parameters, the url parameters must be a number."
    channelInvalidUUID: "Invalid channel ID. It must respect the UUID format."
    channelExists: "Channel already exists."
    wordBlank: "The word cannot be blank."
    wordBlacklist: "This word is not allowed!"
    wordInvalid: "This is not a word, please enter a valid word."
    difficultyRange: "The difficulty must be between 0 and 3."
    hintLimits: "The game must have at least 1 hint and not more than 10."
    hintEmpty: "The hints cannot be empty."
    hintWord: "The hint cannot contain the word."
    hintDuplicate: "The hint cannot be the same."
    invalidUUID: "A valid uuid must be set."
    gameNotFound: "The game cannot be found. Please check if the id is valid."
    gameFileInvalid: "The file is not valid"
    fileNotReadable: "The file cannot be read."
    fileInvalidType: "The file is not a valid type"
    ffaRound: "The number of round must be between 1 and 5 for the free for all game mode."
    soloCount: "The number of players must be one for the game mode solo."
    playerCount: "The number of players must be between 1 and %d."
    playerMax: "You cannot have more than %d players in a game."
    groupSingle: "You already have a group created you cannot create multiple groups."
    groupNotFound: "The group could not be found."
    usernameInvalid: "The username must be between 4 and 12 characters."
    usernameFail: "The user could not be created!"
    usernameExists: "The username already exists. Please choose a different username."
    firstNameInvalid: "Invalid first name or last name, it should not be empty."
    emailInvalid: "Invalid email, it must respect the typical email format."
    passwordInvalid: "Invalid password, it must have a minimum of 8 characters."
    sessionExists: "The user already has an other session tied to this account. Disconnect the other device and wait 5 seconds."
    userNotFound: "The specified user cannot be found."
    userUpdate: "The user could not be updated."
    modifications: "No modifications are found."
    groupOwner: "Only the group owner can kick people out."
    groupMembership: "The user does not belong to a group."
    userLeaveChannel: "User is not in the channel."
    userJoinChannel: "User is already joined to the channel."
    userJoinFirst: "The user needs to join the channel first."
    channelInvalidName: "The channel name is invalid."
    groupIsFull: "The group is full."
    userSingleGroup: "The user can only join one group."
    gameMinimum: "There are only %d game(s). There needs to be a minimum of 10 games before it can start."
    notEnough: "There are not enough users to start the game."
    passwordWrong: "The password is incorrect."
    loginBearer: "The username and the bearer must be set."
    hintInvalid: "Hints are not available for this player. The drawing player needs to be a virtual player."
    hintScore: "You need at least 50 points for a hint."
    hintTime: "There needs to be at least 10 seconds for a hint to be requested."

  botlines:
    startgameAngry1: "I hate this game and I don't like the players in the group. It's an horrible game!"
    startgameAngry2: "Sick of this game, but well I'm going to play since we're confined."
    startgameAngry3: "The game will start, I hope I have fun because it usually doesn't amuse me."

    startgameFunny1: "The game will start, fasten your seatbelts ... Or should I say your keyboards and mouses!"
    startgameFunny2: "Well, I took a new resolution before the game begins: 1080p xD. Enjoy the game ;) !"
    startgameFunny3: "The game will begin. Protect your screen, there are cheaters who wants to raise their stats!"

    startgameMean1: "Before the game begins, I wish you all to lose."
    startgameMean2: "The game will start, I hope your machine will run out of battery before the end!"
    startgameMean3: "I want to warn you before the beginning, I will do everything to make you lose."

    startgameNice1: "Before the game begins, I wish you all good luck!"
    startgameNice2: "The game will start, have fun and may the best win!"
    startgameNice3: "I have a feeling we will have fun! Let's go team!"

    startgameSupportive1: "Before the game begins, remember that the important thing is to participate. Everyone is a winner in the end !"
    startgameSupportive2: "Making mistakes is worth trying! So give everything for this game!"
    startgameSupportive3: "If you feel like giving up, act like that the stats will count in your GPA. It'll give you a boost!"


    endRoundAngry1: "I still don't like this game !!!!"
    endRoundAngry2: "Still confined even after the round. I guess that we'll have to continue ..."
    endRoundAngry3: "The round was horrible. I have a feeling that the next one will be worse."

    endRoundFunny1: "The round is finished. No one is KO yet I suppose?"
    endRoundFunny2: "Before the next round begins, may I go to the bathroom please? I have to empty my RAM."
    endRoundFunny3: "Wow this round was intense. Looks like I've moved mountains. You would be surprised what a virtual bot can do."

    endRoundMean1: "I'm done with this round!  All of you have lost because you're bad."
    endRoundMean2: "Your battery made it through this round, but it will crash the next one? I sure hope so!"
    endRoundMean3: "What was this round guys? I have rarely seen players that bad."

    endRoundNice1: "This round was so cool, wasn't it? Let's go to the next!!"
    endRoundNice2: "Wow what round. I really enjoyed it! We should do it again."
    endRoundNice3: "This round was good, but I'm sure the next one will be even better!"

    endRoundSupportive1: "Everyone did very well in this round. Congratulations and let's go to the next!"
    endRoundSupportive2: "Even if you make mistakes, don't worry. You'll do better the next time."
    endRoundSupportive3: "Come on don't give up! Your GPA is still in play!"


    playerRefAngry1: "I don't know why, but I don't like {}. Your in game duration is []. I can't stand you anymore."
    playerRefAngry2: "It annoys me even confined, {} always plays. Your in game duration is []. I can't stand you anymore."
    playerRefAngry3: "Please someone to kick-out {}? I can smell this player from here. Your in game duration is []. I can't stand you anymore."

    playerRefFunny1: "Hey {}! Always up for a bowling game after the game? Your in game duration is []. You're such a geek :D."
    playerRefFunny2: "Hey {}! Can you lend me your body please? I want to go swimming, but my motherboard isn't water-resistant. Your in game duration is []. You're such a geek :D."
    playerRefFunny3: "I had a joke on {}. But if I say it, then you will say that I'm a racist... Your in game duration is []. You're such a geek :D."

    playerRefMean1: "It's hopeless for {}... Your in game duration is []. Don't like you since second 1."
    playerRefMean2: "I'll hack you {} and retrieve your personal information! Your in game duration is []. Don't like you since second 1."
    playerRefMean3: "Well we can pretend that {} is not playing. It wouldn't change much... Your in game duration is []. Don't like you since second 1."

    playerRefNice1: "I love to play with {}! Your in game duration is []. Glad you like this game !."
    playerRefNice2: "What a blast to play with {}. Your in game duration is []. Glad you like this game !."
    playerRefNice3: "Oh so cool that {} is with us in this game! Your in game duration is []. Glad you like this game !."

    playerRefSupportive1: "Well done {}! Your in game duration is [], keep it up!"
    playerRefSupportive2: "Courage {}! Your in game duration is [], don't give up !"
    playerRefSupportive3: "Your in game duration is []. Keep going {} it will pay off, you'll see!"


    hintRequestAngry1: "A clue ?!! It pisses me off!!"
    hintRequestAngry2: "I hope that you'll find it with this clue. You're making me furious!!!"
    hintRequestAngry3: "The hint is useless! Why even bother asking?!"

    hintRequestFunny1: "Someone ordered a not overcooked clue?"
    hintRequestFunny2: "Oh a clue! Or should I say, January 10? Damn this joke would have worked in french."
    hintRequestFunny3: "A wild hint is approaching. I repeat, a wild hint is approaching."

    hintRequestMean1: "Pfff, you're desperate enough to ask for a hint?"
    hintRequestMean2: "Let's be honest! The hint won't change anything in your incompetence in this game."
    hintRequestMean3: "Someone asked for a clue? Hahaha I was convinced that you were AKF!"

    hintRequestNice1: "There, I hope that this hint will be able to help you."
    hintRequestNice2: "Oh cool, a new clue!"
    hintRequestNice3: "This hint will be useful to you. I'm sure!"

    hintRequestSupportive1: "Good initiative for asking a clue! Keep it up!"
    hintRequestSupportive2: "This hint will help you to find the word! Keep looking!"
    hintRequestSupportive3: "Don't lose hope, that clue might help you!"

    winRatioAngry1: "Hey {} ! You think you're the strongest because you won the most games? You are getting on my nerves !!! (WinRatio of : [])"
    winRatioAngry2: "I don't even know why I'm playing, {} will win anyway!!! (WinRatio of : [])"
    winRatioAngry3: "Can someone kick out {}? He has won enough games! We don't need to help him ANYMORE! (WinRatio of : [])"

    winRatioFunny1: "{} is the MVP of the game! (MVP isn't the concept of Dropbox obviously. I've forgot it anyway). (WinRatio of : [])"
    winRatioFunny2: "This is the story of {} that plays a game and who wins. (My jokes are automated, it's not my fault if they are bad...) (WinRatio of : [])"
    winRatioFunny3: "When I saw the number of games {} won, I lost track. (I mean Threads of course. I'm a bot that lives in a world of semaphores and mutexes.) (WinRatio of : [])"

    winRatioMean1: "Hey {}! It isn't because you won the most games, that you're going to win this one. (WinRatio of : [])"
    winRatioMean2: "The goal is to make lose {}, he won enough games like that !! (WinRatio of : [])"
    winRatioMean3: "{} is the player to shoot! We must bring down his winning ratio. (WinRatio of : [])"

    winRatioNice1: "{} is the best player in the game! (WinRatio of : [])"
    winRatioNice2: "Hey {}! According to your win ratio you're on track to win this game too. (WinRatio of : [])"
    winRatioNice3: "We can all applaud {}! He has the best win ratio among us! (WinRatio of : [])"

    winRatioSupportive1: "Hey {}! You have to let the others win in order to encourage them! (WinRatio of : [])"
    winRatioSupportive2: "Courage {}! Don't give up ! (WinRatio of : [])"
    winRatioSupportive3: "Follow your efforts {} it will pay off, you'll see! (WinRatio of : [])"

    noHintAngry: "You've already asked me all the available hints. What do you want from me? You're pissing me off!!"
    noHintFunny: "Oups, no more clue for you. May I offer you a drink instead? (without alcool, of course. Thanks Tanguy!)"
    noHintMean: "No more hints for this drawing. But to be honnest, I didn't want to give you one."
    noHintNice: "Sorry, I don't have any more hints for this drawing. :("
    noHintSupportive: "I don't have any remaining hints for this drawing. But keep it up, I know you can figure this out !"
    
fr:
  error:
    channelNotFound: "Le canal spécifié n'a pu être trouvé."
    channelInvalidStart: "Paramètres invalides, start doit être celui le plus bas."
    channelInvalidUrl: "Paramètres invalides, les paramètres de l'url doivent être des nombres."
    channelInvalidUUID: "Identifiant de canal invalide. L'identifiant doit respecter le format UUID."
    channelExists: "Le canal existe déjà."
    wordBlank: "Le mot ne peut pas être vide."
    wordBlacklist: "Le mot n'est pas autorisé!"
    wordInvalid: "Ceci n'est pas un mot valide, veuillez entrer un mot valide."
    difficultyRange: "La difficulté doit être entre 0 et 3."
    hintLimits: "Le jeu doit avoir entre 1 et 10 indices maximums."
    hintEmpty: "Les indices ne peuvent être vides."
    hintWord: "L'indice ne peut pas contenir le mot à deviner."
    hintDuplicate: "L'indice ne peut être indentique aux autres."
    invalidUUID: "Un UUID valide doit être utilisé."
    gameNotFound: "Le jeu n'a pas été trouvé. Vérifiez si l'identifiant est valide."
    gameFileInvalid: "Le fichier n'est pas valide."
    fileNotReadable: "Le fichier ne peut être lu."
    fileInvalidType: "Le fichier n'est pas un type valide."
    ffaRound: "Le nombre de tours doit être entre 1 et 5 pour le type mêlée générale."
    soloCount: "Le nombre de joueurs doit être de un pour le mode de jeu solo."
    playerCount: "Le nombre de joueurs doit être entre 1 et %d."
    playerMax: "Vous ne pouvez pas avoir plus de %d joueurs dans une partie."
    groupSingle: "Vous avez déjà créé un groupe. Vous ne pouvez pas avoir plusieurs groupes à votre nom."
    groupNotFound: "Le groupe ne peut pas être trouvé."
    usernameInvalid: "Le nom d'utilisateur est doit être entre 4 et 12 caractères."
    usernameFail: "Le nom d'utilisateur n'a pu être créé!"
    usernameExists: "Le nom d'utilisateur existe déjà. Veuillez en choisir un autre."
    firstNameInvalid: "Prénom ou le nom est invalide. Il ne doit pas être vide."
    emailInvalid: "Courriel invalide, il doit respecter le format d'une adresse courriel."
    passwordInvalid: "Mot de passe invalide. Le mot de passe doit avoir un minimum de 8 caractères."
    sessionExists: "L'utilisateur possède déjà une session. Veuillez déconnecter l'autre session et attendez 5 secondes avant de réessayer."
    userNotFound: "L'utilisateur n'est pas trouvable."
    userUpdate: "L'utilisateur n'a pas été mis à jour."
    modifications: "Les modifications n'ont pas été trouvés."
    groupOwner: "Seulement le propriétaire du groupe peut retirer des joueurs."
    groupMembership: "L'utilisateur n'appartient pas à un groupe."
    userLeaveChannel: "L'utilisateur n'est pas dans le canal."
    userJoinChannel: "L'utilisateur est déjà présent dans le canal."
    userJoinFirst: "L'utilisateur doit rejoindre le canal avant."
    channelInvalidName: "Le nom du canal est invalide."
    groupIsFull: "Le groupe est plein."
    userSingleGroup: "L'utilisateur ne peut joindre qu'un seul groupe."
    gameMinimum: "Il n'a seulement que %d jeu(x). Un minimum de 10 jeux doivent être présents pour débuter la partie."
    notEnough: "Il n'a pas assez d'utilisateurs pour démarrer la partie."
    passwordWrong: "Le mot de passe est erroné."
    loginBearer: "Le nom d'utilisateur et le bearer doivent être dans la requête."
    hintInvalid: "Les indices ne sont pas disponible pour ce joueur. Ils ne sont disponibles que pour les joueurs virtuels."
    hintScore: "Vous avez besoin de 50 pts pour demander un indice. Vous n'avez pas assez de points."
    hintTime: "Il doit vous rester au minimum 10 secondes pour demander un indice."

  botlines:
    startgameAngry1: "J'aime pas ce jeu et je n'aime pas les joueurs du groupe. Mauvais jeu à tous !"
    startgameAngry2: "Ras le bol de ce jeu, mais bon je vais quand même jouer vu qu'on est confiné."
    startgameAngry3: "La partie va commencer, j'espère que je vais prendre du plaisir car d'habitude je ne m'amuse pas."
    
    startgameFunny1: "La partie va commencer, accrochez vos ceintures ... Ou devrais-je dire à votre clavier et souris !"
    startgameFunny2: "Bon j'ai pris une nouvelle résolution pour la partie qui va commencer : Du 1080p xD. Bonne partie ;) !"
    startgameFunny3: "La partie va commencer. Protège ton écran des tricheurs qui veulent faire monter leur stats !"
    
    startgameMean1: "Avant que la partie commence, je vous souhaite tous de perdre."
    startgameMean2: "La partie va commencer, j'espère que ton appareil va manquer de batterie avant la fin !"
    startgameMean3: "Je tiens à vous prévenir avant le commencement, je vais tout faire pour vous faire perdre."
    
    startgameNice1: "Avant que la partie commence, je vous souhaite à tous bonne chance !"
    startgameNice2: "La partie va commencer, prenez du plaisir et que le meilleur gagne !"
    startgameNice3: "Je sens qu'on va bien s'amuser ! C'est parti !"
    
    startgameSupportive1: "Avant que la partie commence, n'oubliez pas que l'important c'est de participer. Tout le monde sort gagnant de la partie !"
    startgameSupportive2: "Faire des erreurs, c'est la preuve d’avoir essayé ! Alors donne tout pour cette partie !"
    startgameSupportive3: "Si tu as envie d'abandonner, dis toi que les stats vont compter dans ton GPA. Ça va te donner un boost !"
    
    
    endRoundAngry1: "J'aime toujours pas ce jeu !!!!"
    endRoundAngry2: "Toujours confiné, même après le round. I guess qu'il faut qu'on continue..."
    endRoundAngry3: "Le round était nul. Mais je sens que le prochain va être pire."
    
    endRoundFunny1: "Le round est fini. Personne n'est encore KO j'espère ?"
    endRoundFunny2: "Avant que le prochain round commence, je peux aller aux toilettes svp ? Je dois vider ma RAM."
    endRoundFunny3: "Wow c'était physique ce round. On dirait que j'ai soulevé des montagnes. Et oui, vous serez surpris de ce qu'un bot virtuel peut faire."
    
    endRoundMean1: "Le round est fini, pour moi vous avez tous perdu car vous êtes tous nuls."
    endRoundMean2: "Ta batterie a tenu pour ce round là, mais elle va lâcher au prochain (du moins je le souhaite)."
    endRoundMean3: "C'était quoi ce round guys ? J'ai rarement vu des joueurs aussi mauvais."
    
    endRoundNice1: "Ce round était trop cool. On passe au prochain ?"
    endRoundNice2: "Wow quel round. Je prends vraiment du plaisir."
    endRoundNice3: "Ce round était bon, mais je suis sûr que le prochain sera encore mieux !"
    
    endRoundSupportive1: "Tout le monde s'est très bien débrouillé dans ce round. Bravo et passons au prochain !"
    endRoundSupportive2: "Même si tu as fait des erreurs, ne t'inquiète pas. Tu feras mieux au prochain."
    endRoundSupportive3: "Aller on lâche rien ! Ton GPA est toujours en jeu !"
    
    
    playerRefAngry1: "Je sais pas pourquoi mais {}, je ne l'aime pas. Ton temps total de jeu est de []. Je ne te supporte plus."
    playerRefAngry2: "Ça m'énerve même confiné {} joue toujours. Ton temps total de jeu est de []. Je ne te supporte plus."
    playerRefAngry3: "Bon on peut kick-out {} svp ? Il m'enerve. Ton temps total de jeu est de []. Je ne te supporte plus."
    
    playerRefFunny1: "Hey {} ! Toujours partant pour faire un bowling après la partie ? Ton temps total de jeu est de []. Tu es un geek :D."
    playerRefFunny2: "{} peux-tu me preter ton corps stp ? Je veux aller à la piscine mais ma carte mère est pas resistante à l'eau. Ton temps total de jeu est de []. Tu es un geek :D."
    playerRefFunny3: "J'avais une blague sur {} mais après on va dire que je suis raciste ... Ton temps total de jeu est de []. Tu es un geek :D."
    
    playerRefMean1: "Hey {} ! Tu es nul ! Ton temps total de jeu est de []. Je ne t'aime pas depuis la première seconde."
    playerRefMean2: "Je vais hack {} et récupèrer tes informations personnelles ! Ton temps total de jeu est de []. Je ne t'aime pas depuis la première seconde."
    playerRefMean3: "Bon on peut faire comme si que {} jouait pas. Ça serait pareil ... Ton temps total de jeu est de []. Je ne t'aime pas depuis la première seconde."
    
    playerRefNice1: "J'adore faire des parties avec {} ! Ton temps total de jeu est de []. Content que tu apprécis le jeu !"
    playerRefNice2: "Quel plaisir de jouer en présence de {}. Ton temps total de jeu est de []. Content que tu apprécis le jeu !"
    playerRefNice3: "Oh trop cool que {} soit avec nous dans cette partie ! Ton temps total de jeu est de []. Content que tu apprécis le jeu !"
    
    playerRefSupportive1: "Ton temps total de jeu est de []. Bien joué {}, continue comme ça !"
    playerRefSupportive2: "Ton temps total de jeu est de []. Courage {} ! Ne lâche rien !"
    playerRefSupportive3: "Ton temps total de jeu est de []. Poursuis tes efforts {} ça va payer, tu vas voir !"
    
    
    hintRequestAngry1: "Un indice ?!! Ça m'énerve !!"
    hintRequestAngry2: "J'espère qu'avec cet indice tu vas trouver car la je suis à deux doigts de m'énerver !!"
    hintRequestAngry3: "L'indice est nul ! Ça m'énerve !!"
    
    hintRequestFunny1: "Quelqu'un a commandé un indice pas trop cuit ?"
    hintRequestFunny2: "Oh un indice ! Ou devrais-je dire 1 10 ?"
    hintRequestFunny3: "Indice sauvage en approche. Je répète, indice sauvage en approche."
    
    hintRequestMean1: "Tssss vous êtes assez désespéré pour demander un indice ?"
    hintRequestMean2: "Soyons honnête, l'indice ne va rien changer quand à votre incompétence pour trouver le mot."
    hintRequestMean3: "Quelqu'un a demandé un indice ? Ahaha j'étais persuadé que vous n'étiez pas à la hauteur pour ce jeu."
    
    hintRequestNice1: "J'espère que cet indice va pouvoir t'aider."
    hintRequestNice2: "Ah cool, un nouvel indice !"
    hintRequestNice3: "Cet indice va sûrement t'être utile !"
    
    hintRequestSupportive1: "Bonne initiative d'avoir demandé un indice ! Continue comme ça !"
    hintRequestSupportive2: "Cet indice va t'aider pour trouver le mot ! Continue de chercher !"
    hintRequestSupportive3: "Ne perd pas espoir, voilà un indice pour t'aider !"
    
    winRatioAngry1: "Hey {} tu te crois plus fort parce que tu as gagné le plus de parties ? Tu m'énerves !!! (Ratio de victoire de : [])"
    winRatioAngry2: "Je sais même pas pourquoi je joue, c'est encore {} qui va gagner !!! (Ratio de victoire de : [])"
    winRatioAngry3: "Bon on peut kick-out {} svp ? Il a déjà gagné assez de parties comme ça !! (Ratio de victoire de : [])"
    
    winRatioFunny1: "{} c'est le MVP de la partie ! (MVP c'est pas le concept de dropbox evidemment) (Ratio de victoire de : [])"
    winRatioFunny2: "C'est l'histoire de {} qui joue à une partie. Et qui l'a gagne. (mes blagues sont automatisées, c'est pas ma faute ...) (Ratio de victoire de : [])"
    winRatioFunny3: "Quand je vois le nombre de partie que {} a gagné, j'en perds le fil. (par fil j'entends thread evidemment. On oublie pas que je suis qu'un simple joueur virtuel) (Ratio de victoire de : [])"
    
    winRatioMean1: "Hey {} C'est pas parce que ta gagné le plus de parties ici que tu vas gagner celle la. (Ratio de victoire de : [])"
    winRatioMean2: "Le but c'est de faire perdre {}, il a gagné suffisamment de parties comme ça !! (Ratio de victoire de : [])"
    winRatioMean3: "{} c'est le joueur à abattre ! On doit faire baisser son ratio de victoires. (Ratio de victoire de : [])"
    
    winRatioNice1: "{} c'est le meilleur joueur de la partie ! (Ratio de victoire de : [])"
    winRatioNice2: "{} je pense que t'es bien parti pour gagner cette partie aussi. (Ratio de victoire de : [])"
    winRatioNice3: "On peut tous applaudir {} ! C'est lui qui a gagné le plus de parties parmi nous ! (Ratio de victoire de : [])"
    
    winRatioSupportive1: "{} il faut que tu laisses les autres gagner afin de les encourager ! (Ratio de victoire de : [])"
    winRatioSupportive2: "Courage {} ! Ne lâche rien ! (Ratio de victoire de : [])"
    winRatioSupportive3: "Poursuis tes efforts {} ça va payer, tu vas voir ! (Ratio de victoire de : [])"

    noHintAngry: "Tu as déjà demandé tous les indices disponibles. Que veux-tu de plus ?! Tu m'énerves"
    noHintFunny: "Oups, j'ai plus d'indices pour toi. Puis-je t'offrir un verre à la place? (sans alcool bien sûr)"
    noHintMean: "Il n'y a plus d'indices pour ce dessin. Mais pour être honnête, je voulais pas t'en donner un."
    noHintNice: "Désolé, je n'ai plus d'indices pour ce dessin :("
    noHintSupportive: "Je n'ai plus d'indices pour ce dessin. Mais tiens bon, je sais que tu peux le trouver !"
`)
