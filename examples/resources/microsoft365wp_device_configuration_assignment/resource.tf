resource "microsoft365wp_device_configuration_assignment" "example" {
  device_configuration_id = "11111111-1111-4111-8111-111111111111"
  intent                  = "apply"
  target = {
    group = {
      group_id = "22222222-2222-4222-8222-222222222222"
    }
  }
}
