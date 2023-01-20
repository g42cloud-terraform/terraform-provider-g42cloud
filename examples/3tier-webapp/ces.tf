resource "g42cloud_smn_topic" "topic" {
  name         = "smn-rds-pg"
  display_name = "smn topic for pg"
}

# CES alarm rule for CPU Usage > 80%
resource "g42cloud_ces_alarmrule" "alarmrule_cpu" {
  alarm_name           = "rule-pg-cpu"
  alarm_action_enabled = true

  metric {
    namespace   = "SYS.RDS"
    metric_name = "rds001_cpu_util"

    dimensions {
      name  = "postgresql_instance_id"
      value = g42cloud_rds_instance.rds.id
    }
  }

  condition  {
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 80
    unit                = "%"
    count               = 1
  }

  alarm_actions {
    type              = "notification"
    notification_list = [
      g42cloud_smn_topic.topic.topic_urn
    ]
  }
}

# CES alarm rule for MEM Usage > 85%
resource "g42cloud_ces_alarmrule" "alarmrule_mem" {
  alarm_name           = "rule-pg-mem"
  alarm_action_enabled = true

  metric {
    namespace   = "SYS.RDS"
    metric_name = "rds001_cpu_util"

    dimensions {
      name  = "postgresql_instance_id"
      value = g42cloud_rds_instance.rds.id
    }
  }

  condition  {
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 85
    unit                = "%"
    count               = 1
  }

  alarm_actions {
    type              = "notification"
    notification_list = [
      g42cloud_smn_topic.topic.topic_urn
    ]
  }
}

# CES alarm rule for DISK Usage > 80%
resource "g42cloud_ces_alarmrule" "alarmrule_disk" {
  alarm_name           = "rule-pg-disk"
  alarm_action_enabled = true

  metric {
    namespace   = "SYS.RDS"
    metric_name = "rds001_cpu_util"

    dimensions {
      name  = "postgresql_instance_id"
      value = g42cloud_rds_instance.rds.id
    }
  }

  condition  {
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 80
    unit                = "%"
    count               = 1
  }

  alarm_actions {
    type              = "notification"
    notification_list = [
      g42cloud_smn_topic.topic.topic_urn
    ]
  }
}
