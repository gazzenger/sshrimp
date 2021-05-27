
data "aws_iam_policy_document" "sshrimp_ca_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
  statement {
    actions = ["sts:AssumeRoleWithWebIdentity"]

    principals {
      type        = "Federated"
      identifiers = ["${var.webidentity_principal_identifiers}"]
    }

    condition {
      test     = "StringEquals"
      variable = "${var.webidentity_provider_url}"
      values = ["${var.webidentity_client_id}"]
    }
  }
}

data "aws_iam_policy_document" "sshrimp_ca" {
  statement {
    actions = [
      "kms:Sign",
      "kms:GetPublicKey"
    ]
    resources = [
      "${aws_kms_key.sshrimp_ca_private_key.arn}",
    ]
  }

  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
      "lambda:InvokeFunction",
    ]

    resources = [
      "*",
    ]
  }
}

resource "aws_iam_role_policy" "sshrimp_ca" {
  name   = "sshrimp-ca-${data.aws_region.current.name}"
  role   = aws_iam_role.sshrimp_ca.id
  policy = data.aws_iam_policy_document.sshrimp_ca.json
}

resource "aws_iam_role" "sshrimp_ca" {
  name               = "sshrimp-ca-${data.aws_region.current.name}"
  assume_role_policy = data.aws_iam_policy_document.sshrimp_ca_assume_role.json
}
