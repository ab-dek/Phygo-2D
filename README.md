# Phygo2D

A simple, lightweight, and efficient 2D rigid body physics engine written in Go. 

Phygo2D is designed to handle basic 2D physics simulations, including collision detection, resolution, gravity, and friction. It is decoupled from any rendering engine, allowing you to use it with [Raylib](https://github.com/gen2brain/raylib-go) or [Ebiten](https://github.com/hajimehoshi/ebiten).

## Features

* **Shape Support:** currently only Circles and Rectangles 
* **Collision Detection:**
    * Uses Separating Axis Theorem for accurate Polygon-Polygon and Polygon-Circle detection.
    * AABB broad-phase optimization.
* **Physical Properties:**
    * Mass, Density, and Restitution (Bounciness).
    * Static and Dynamic Friction.
    * Gravity and Force application.

## Installation

To use Phygo2D in your Go project, run:
```bash
go get https://github.com/ab-dek/phygo2d
```

## Usage
Check the examples [here](./example/)