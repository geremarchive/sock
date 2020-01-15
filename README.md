<h1 align="center">sock 🔒</h1>

<p align="center">a simple terminal locker</p>
<img align="right" src="media/sock.png">

## About

```sock```, or **S**imple L**ock** is a simple and customizable terminal locker. ```sock``` allows you to lock one or multiple terminals and sync them. It matches the password the user types in with the encypted root password stored in ```/etc/shadow``` 

## Usage

(must be run a root)

```
Usage: sock [OPTION]
A simple terminal locker

--help, -h: Display this information
--message=[STRING], -m: Set the message
--color=[HEX], -c: Set the color of the message
--bg=[HEX], -b: Set the color of the background
--og=[HEX], -o: Set the original color (only needed if you are using the '-e' option)
--escape, -e: Set the background color using escape seqences (check if your terminal supports this)
--center, -C: Center the text
--bold, -B: Use bold text
--all, -a: Lock all terminals
--check, -k: Checks if all terminals are locked, if they are, it locks the terminal
```

## Setup

Add ```sock -k``` to your ```.shellrc``` so newly opened terminals can check if all terminals are locked

## Dependencies

```
github.com/gdamore/tcell
```

```
github.com/spf13/pflag
```

```
openssl
```
