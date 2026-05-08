using Godot;
using System;

public partial class Player : CharacterBody3D
{
	public const float Speed = 5.0f;
	public const float JumpVelocity = 4.5f;
	public const float MouseSensitivity = 0.002f;


	private Camera3D _camera;

	public override void _Ready()
	{
		_camera = GetNode<Camera3D>("Camera3D");
		Input.MouseMode = Input.MouseModeEnum.Captured;
	}

	public override void _Input(InputEvent @event)
	{

		if(Input.MouseMode != Input.MouseModeEnum.Captured) return;

		if (@event is InputEventMouseMotion motion)
		{
			// left/right (whole body)
			RotateY(-motion.Relative.X * MouseSensitivity);

			// up/down (camera tilt)
			_camera.RotateX(-motion.Relative.Y * MouseSensitivity);
			_camera.Rotation = new Vector3(
					Mathf.Clamp(_camera.Rotation.X, -1.2f, 1.2f),
					_camera.Rotation.Y,
					_camera.Rotation.Z
			);
		}

	}

	public override void _PhysicsProcess(double delta)
	{
		Vector3 velocity = Velocity;

		if(Input.MouseMode != Input.MouseModeEnum.Captured){
			Velocity = Vector3.Zero;
			MoveAndSlide();
			return;
		}

		// Add the gravity.
		if (!IsOnFloor())
		{
			velocity += GetGravity() * (float)delta;
		}

		// Handle Jump.
		if (Input.IsActionJustPressed("ui_accept") && IsOnFloor())
		{
			velocity.Y = JumpVelocity;
		}

		// Get the input direction and handle the movement/deceleration.
		// As good practice, you should replace UI actions with custom gameplay actions.
		Vector2 inputDir = Input.GetVector("ui_left", "ui_right", "ui_up", "ui_down");
		Vector3 direction = (Transform.Basis * new Vector3(inputDir.X, 0, inputDir.Y)).Normalized();
		if (direction != Vector3.Zero)
		{
			velocity.X = direction.X * Speed;
			velocity.Z = direction.Z * Speed;
		}
		else
		{
			velocity.X = Mathf.MoveToward(Velocity.X, 0, Speed);
			velocity.Z = Mathf.MoveToward(Velocity.Z, 0, Speed);
		}

		Velocity = velocity;
		MoveAndSlide();
	}
}
