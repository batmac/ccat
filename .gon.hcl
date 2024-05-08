source    = ["./ccat"]
bundle_id = "com.sbz.ccat"

sign {
  application_identity = "Developer ID Application"
}

# for stapling
dmg {
  output_path = "ccat.dmg"
  volume_name = "ccat"
}

zip {
  output_path = "ccat.zip"
}
