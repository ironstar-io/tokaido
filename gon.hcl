{
    "source" : ["./dist/tok-macos"],
    "bundle_id" : "io.tokaido",
    "apple_id": {
        "password":  "@env:AC_PASSWORD"
    },
    "sign" :{
        "application_identity" : "Developer ID Application: Ironstar Hosting Services Pty Ltd"
    },
    "dmg" :{
        "output_path":  "./dist/tok-macos.dmg",
        "volume_name":  "Tokaido"
    },
    "zip" :{
        "output_path" : "./dist/tok-macos.zip"
    }
}
