using Godot;
using System;

public partial class Computer : StaticBody3D
{
	public bool _playerNearby = false;
	public bool _terminalActive = false;
	private LineEdit _input;
	private SubViewport _screenViewport;

	// Called when the node enters the scene tree for the first time.
	public override void _Ready()
	{

		_input = GetNode<LineEdit>("ScreenViewport/TerminalUI/VBoxContainer/InputRow/Input");
		_screenViewport = GetNode<SubViewport>("ScreenViewport");


		var zone = GetNode<Area3D>("InteractionZone");
		var output = GetNode<RichTextLabel>("ScreenViewport/TerminalUI/VBoxContainer/Output");
		output.GetVScrollBar().Modulate = new Color(0, 0, 0, 0);

		zone.BodyEntered += OnBodyEntered;
		zone.BodyExited += OnBodyExited;

		_input.TextSubmitted += OnCommandSubmitted;

		var result = AxiomBridge.Init("../engine/cmd/axiom/initial_config.ax");
		if (result != "")
			GD.Print("AxiomInit failed: " + result);

	}

	private void ScrollToBottom()
	{
	}

	private void OnBodyEntered(Node3D body)
	{
		if (body is CharacterBody3D)
		{
			_playerNearby = true;
		}
	}

	private void OnBodyExited(Node3D body)
	{
		if (body is CharacterBody3D)
		{
			_playerNearby = false;
			Deactivate();
		}
	}

	public override void _UnhandledInput(InputEvent @event)
	{

		if (_terminalActive)
		{
			if (@event is InputEventKey key && key.Pressed)
			{
				if (key.Keycode == Key.Escape)
				{
					Deactivate();
					return;
				}

				if (key.Keycode == Key.Enter || key.Keycode == Key.KpEnter)
				{
					OnCommandSubmitted(_input.Text);
					_input.Clear();
					GetViewport().SetInputAsHandled();
					return;
				}
			}

			_screenViewport.PushInput(@event);
			GetViewport().SetInputAsHandled();
			return;
		}

		if (_playerNearby && @event is InputEventKey e && e.Pressed && e.Keycode == Key.E)
		{
			Activate();
		}


	}

	public void Activate()
	{
		_terminalActive = true;
		Input.MouseMode = Input.MouseModeEnum.Visible;
		_input.GrabFocus();
	}


	public void Deactivate()
	{
		_terminalActive = false;
		Input.MouseMode = Input.MouseModeEnum.Captured;
		_input.ReleaseFocus();
	}

	private void OnCommandSubmitted(string command)
	{
		var output = GetNode<RichTextLabel>("ScreenViewport/TerminalUI/VBoxContainer/Output");
		var result = AxiomBridge.Execute(command);

		output.AppendText($"\n> {command}\n{result}\n");
		CallDeferred(nameof(ScrollToBottom));
		_input.Clear();
	}

	// Called every frame. 'delta' is the elapsed time since the previous frame.
	public override void _Process(double delta)
	{
	}
}
