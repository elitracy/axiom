using System.Runtime.InteropServices;
using System.Reflection;
using System;
using Godot;

public static class AxiomBridge
{

    static AxiomBridge(){
        NativeLibrary.SetDllImportResolver(typeof(AxiomBridge).Assembly, Resolver);
    }

    public static IntPtr Resolver(string name, Assembly assembly, DllImportSearchPath? path) {
        if (name != "axiom") return IntPtr.Zero;
        var libPath = System.IO.Path.Combine(
                ProjectSettings.GlobalizePath("res://"),
                "libaxiom.dylib"
                );

        return NativeLibrary.Load(libPath);
    }

    [DllImport("axiom")]
    private static extern IntPtr AxiomInit(string configPath);

    [DllImport("axiom")]
    private static extern IntPtr AxiomExecute(string input);

    public static string Init(string configPath)
    {
        var ptr = AxiomInit(configPath);
        return Marshal.PtrToStringAnsi(ptr) ?? "";
    }

    public static string Execute(string input)
    {
        var ptr = AxiomExecute(input);
        return Marshal.PtrToStringAnsi(ptr) ?? "";
    }
}
