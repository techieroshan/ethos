package templates

// Email Verification Template - 95%+ Email Client Compatibility
const emailVerificationTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify Your Email - Ethos Platform</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }
        .container { max-width: 600px; margin: 0 auto; background-color: #ffffff; }
        .header { background-color: #667eea; padding: 40px 30px; text-align: center; }
        .logo { width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px; }
        .logo div { width: 40px; height: 25px; border-bottom: 3px solid #667eea; border-radius: 2px; margin: 0 auto; }
        .content { padding: 40px 30px; }
        .button { display: inline-block; background-color: #667eea; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600; margin: 20px 0; }
        .footer { background-color: #f9fafb; padding: 20px 30px; text-align: center; font-size: 14px; color: #6b7280; }
    </style>
</head>
<body>
    <!-- Outlook-safe table-based layout for 95%+ compatibility -->
    <table class="container" width="100%" border="0" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: 0 auto; background-color: #ffffff;">
        <tr>
            <td class="header" style="background-color: #667eea; padding: 40px 30px; text-align: center;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td style="text-align: center;">
                            <div class="logo" style="width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px;">
                                <div style="width: 40px; height: 25px; border-bottom: 3px solid #667eea; border-radius: 2px; margin: 0 auto;"></div>
                            </div>
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Verify Your Email</h1>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>

        <tr>
            <td class="content" style="padding: 40px 30px;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td>
                            <h2 style="color: #1f2937; margin-bottom: 20px; font-size: 24px;">Welcome to Ethos, {{.FirstName}}!</h2>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 20px; font-size: 16px;">
                                Thank you for joining Ethos Platform. To complete your registration and ensure the security of your account, please verify your email address.
                            </p>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                Click the button below to verify your email address and activate your account:
                            </p>

                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="text-align: center; margin: 30px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="{{.VerifyURL}}" style="display: inline-block; background-color: #667eea; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600;">Verify Email Address</a>
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #6b7280; font-size: 14px; margin-bottom: 20px;">
                                If the button doesn't work, you can copy and paste this link into your browser:
                            </p>

                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #f3f4f6; border-radius: 4px; margin: 10px 0;">
                                <tr>
                                    <td style="padding: 10px; font-size: 14px; color: #6b7280; word-break: break-all;">
                                        {{.VerifyURL}}
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #4b5563; line-height: 1.6; font-size: 16px;">
                                This verification link will expire in 24 hours for security reasons. If you didn't create an account with Ethos, please ignore this email.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td style="padding: 20px; text-align: center; font-size: 12px; color: #666; background-color: #f9fafb;">
                ` + canSpamFooter + `
            </td>
        </tr>
    </table>
</body>
</html>`

// Password Reset Template - 95%+ Email Client Compatibility
const passwordResetTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Your Password - Ethos Platform</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }
        .container { max-width: 600px; margin: 0 auto; background-color: #ffffff; }
        .header { background-color: #f59e0b; padding: 40px 30px; text-align: center; }
        .logo { width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px; }
        .logo div { width: 40px; height: 25px; border-bottom: 3px solid #f59e0b; border-radius: 2px; margin: 0 auto; }
        .content { padding: 40px 30px; }
        .button { display: inline-block; background-color: #f59e0b; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600; margin: 20px 0; }
        .warning { background-color: #fef3c7; border: 1px solid #f59e0b; border-radius: 8px; padding: 20px; margin: 20px 0; }
    </style>
</head>
<body>
    <!-- Outlook-safe table-based layout for 95%+ compatibility -->
    <table class="container" width="100%" border="0" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: 0 auto; background-color: #ffffff;">
        <tr>
            <td class="header" style="background-color: #f59e0b; padding: 40px 30px; text-align: center;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td style="text-align: center;">
                            <div class="logo" style="width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px;">
                                <div style="width: 40px; height: 25px; border-bottom: 3px solid #f59e0b; border-radius: 2px; margin: 0 auto;"></div>
                            </div>
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Reset Your Password</h1>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>

        <tr>
            <td class="content" style="padding: 40px 30px;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td>
                            <h2 style="color: #1f2937; margin-bottom: 20px; font-size: 24px;">Password Reset Request</h2>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 20px; font-size: 16px;">
                                Hi {{.Name}}, we received a request to reset your password for your Ethos Platform account.
                            </p>

                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fef3c7; border: 1px solid #f59e0b; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <p style="margin: 0; color: #92400e; font-weight: 600;">Security Notice</p>
                                        <p style="margin: 5px 0 0 0; color: #92400e;">
                                            If you didn't request this password reset, please ignore this email. Your password will remain unchanged.
                                        </p>
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                Click the button below to reset your password. This link will expire in 1 hour for your security.
                            </p>

                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="text-align: center; margin: 30px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="{{.ResetURL}}" style="display: inline-block; background-color: #f59e0b; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600;">Reset Password</a>
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #6b7280; font-size: 14px; margin-bottom: 20px;">
                                If the button doesn't work, copy and paste this link into your browser:
                            </p>

                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #f3f4f6; border-radius: 4px; margin: 10px 0;">
                                <tr>
                                    <td style="padding: 10px; font-size: 14px; color: #6b7280; word-break: break-all;">
                                        {{.ResetURL}}
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #4b5563; line-height: 1.6; font-size: 16px;">
                                For your security, this password reset link can only be used once and will expire in 1 hour.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td style="padding: 20px; text-align: center; font-size: 12px; color: #666; background-color: #f9fafb;">
                ` + canSpamFooter + `
            </td>
        </tr>
    </table>
</body>
</html>`

// Welcome Standard User Template
const welcomeStandardUserTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to Ethos - Your Account is Ready</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }
        .container { max-width: 600px; margin: 0 auto; background-color: #ffffff; }
        .header { background-color: #10b981; padding: 40px 30px; text-align: center; }
        .logo { width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px; }
        .logo div { width: 40px; height: 25px; border-bottom: 3px solid #10b981; border-radius: 2px; margin: 0 auto; }
        .content { padding: 40px 30px; }
        .button { display: inline-block; background-color: #10b981; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600; margin: 20px 0; }
        .feature { background-color: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 8px; padding: 20px; margin: 15px 0; }
        .feature h3 { margin: 0 0 10px 0; color: #065f46; }
        .feature p { margin: 0; color: #047857; }
    </style>
</head>
<body>
    <!-- Outlook-safe table-based layout for 95%+ compatibility -->
    <table class="container" width="100%" border="0" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: 0 auto; background-color: #ffffff;">
        <tr>
            <td class="header" style="background-color: #10b981; padding: 40px 30px; text-align: center;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td style="text-align: center;">
                            <div class="logo" style="width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px;">
                                <div style="width: 40px; height: 25px; border-bottom: 3px solid #10b981; border-radius: 2px; margin: 0 auto;"></div>
                            </div>
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Welcome to Ethos!</h1>
                            <p style="color: #ecfdf5; margin: 10px 0 0 0; opacity: 0.9;">Your journey to better feedback starts here</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>

        <tr>
            <td class="content" style="padding: 40px 30px;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td>
                            <h2 style="color: #1f2937; margin-bottom: 20px; font-size: 24px;">Hello {{.Name}}, Welcome Aboard! 🎉</h2>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                Your Ethos account has been successfully created and verified. You're now part of a community dedicated to meaningful feedback and continuous improvement.
                            </p>

                            <!-- Feature 1 -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #065f46; font-size: 18px;">✨ Share Your Voice</h3>
                                        <p style="margin: 0; color: #047857; font-size: 16px;">Submit feedback, share ideas, and contribute to discussions that matter to you.</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- Feature 2 -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #065f46; font-size: 18px;">🔍 Discover Insights</h3>
                                        <p style="margin: 0; color: #047857; font-size: 16px;">Explore feedback from others, learn from different perspectives, and stay informed.</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- Feature 3 -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #065f46; font-size: 18px;">🔔 Stay Connected</h3>
                                        <p style="margin: 0; color: #047857; font-size: 16px;">Receive notifications about topics that interest you and engage with the community.</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- CTA Button -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="text-align: center; margin: 40px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="{{.DashboardURL}}" style="display: inline-block; background-color: #10b981; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600; font-size: 16px;">Get Started →</a>
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #4b5563; line-height: 1.6; font-size: 16px;">
                                Need help getting started? Check out our <a href="#" style="color: #10b981; text-decoration: none;">help center</a> or reach out to our support team.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td style="padding: 20px; text-align: center; font-size: 12px; color: #666; background-color: #f9fafb;">
                ` + canSpamFooter + `
            </td>
        </tr>
    </table>
</body>
</html>`

// Welcome Organization Admin Template
const welcomeOrgAdminTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to Ethos - Organization Admin Access</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }
        .container { max-width: 600px; margin: 0 auto; background-color: #ffffff; }
        .header { background-color: #8b5cf6; padding: 40px 30px; text-align: center; }
        .logo { width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px; }
        .logo div { width: 40px; height: 25px; border-bottom: 3px solid #8b5cf6; border-radius: 2px; margin: 0 auto; }
        .content { padding: 40px 30px; }
        .button { display: inline-block; background-color: #8b5cf6; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600; margin: 20px 0; }
        .admin-feature { background-color: #faf5ff; border: 1px solid #d8b4fe; border-radius: 8px; padding: 20px; margin: 15px 0; }
        .admin-feature h3 { margin: 0 0 10px 0; color: #6b21a8; }
        .admin-feature p { margin: 0; color: #7c3aed; }
    </style>
</head>
<body>
    <!-- Outlook-safe table-based layout for 95%+ compatibility -->
    <table class="container" width="100%" border="0" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: 0 auto; background-color: #ffffff;">
        <tr>
            <td class="header" style="background-color: #8b5cf6; padding: 40px 30px; text-align: center;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td style="text-align: center;">
                            <div class="logo" style="width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px;">
                                <div style="width: 40px; height: 25px; border-bottom: 3px solid #8b5cf6; border-radius: 2px; margin: 0 auto;"></div>
                            </div>
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Organization Admin Access</h1>
                            <p style="color: #f3e8ff; margin: 10px 0 0 0; opacity: 0.9;">Manage your team's Ethos experience</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>

        <tr>
            <td class="content" style="padding: 40px 30px;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td>
                            <h2 style="color: #1f2937; margin-bottom: 20px; font-size: 24px;">Welcome, {{.Name}}! 🏢</h2>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                You've been granted Organization Administrator privileges on the Ethos Platform. This role gives you the tools to manage your organization's settings, team members, and feedback workflows.
                            </p>

                            <!-- Admin Feature 1 -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #faf5ff; border: 1px solid #d8b4fe; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #6b21a8; font-size: 18px;">👥 Team Management</h3>
                                        <p style="margin: 0; color: #7c3aed; font-size: 16px;">Invite team members, manage roles, and oversee user permissions within your organization.</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- Admin Feature 2 -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #faf5ff; border: 1px solid #d8b4fe; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #6b21a8; font-size: 18px;">⚙️ Organization Settings</h3>
                                        <p style="margin: 0; color: #7c3aed; font-size: 16px;">Configure notification preferences, branding, and workflow settings for your team.</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- Admin Feature 3 -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #faf5ff; border: 1px solid #d8b4fe; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #6b21a8; font-size: 18px;">📊 Analytics & Reporting</h3>
                                        <p style="margin: 0; color: #7c3aed; font-size: 16px;">Access detailed analytics about your organization's engagement and feedback trends.</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- Admin Feature 4 -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #faf5ff; border: 1px solid #d8b4fe; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #6b21a8; font-size: 18px;">🛡️ Compliance & Security</h3>
                                        <p style="margin: 0; color: #7c3aed; font-size: 16px;">Monitor compliance requirements and maintain security standards for your organization.</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- CTA Button -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="text-align: center; margin: 40px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="{{.DashboardURL}}" style="display: inline-block; background-color: #8b5cf6; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600;">Access Admin Dashboard →</a>
                                    </td>
                                </tr>
                            </table>

                            <!-- Security Notice -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fef3c7; border: 1px solid #f59e0b; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #92400e; font-size: 18px;">Important Security Notice</h3>
                                        <p style="margin: 0; color: #92400e; font-size: 16px;">
                                            As an organization administrator, you have access to sensitive information. Please ensure your account remains secure and report any suspicious activity immediately.
                                        </p>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td style="padding: 20px; text-align: center; font-size: 12px; color: #666; background-color: #f9fafb;">
                ` + canSpamFooter + `
            </td>
        </tr>
    </table>
</body>
</html>`

// Welcome Platform Admin Template
const welcomePlatformAdminTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to Ethos - Platform Admin Access</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }
    </style>
</head>
<body>
    <table width="100%" border="0" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: 0 auto; background-color: #ffffff;">
        <tr>
            <td style="background-color: #ef4444; padding: 40px 30px; text-align: center;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td style="text-align: center;">
                            <div style="width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px;">
                                <div style="width: 40px; height: 25px; border-bottom: 3px solid #ef4444; border-radius: 2px; margin: 0 auto;"></div>
                            </div>
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Platform Administrator</h1>
                            <p style="color: #fee2e2; margin: 10px 0 0 0; opacity: 0.9;">Global platform oversight and management</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>

        <tr>
            <td style="padding: 40px 30px;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td>
                            <h2 style="color: #1f2937; margin-bottom: 20px; font-size: 24px;">Welcome, {{.Name}}! 🌐</h2>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                You have been granted Platform Administrator privileges. This is the highest level of access on the Ethos Platform, giving you comprehensive control over all system operations, user management, and platform-wide settings.
                            </p>

                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #991b1b; font-size: 18px;">🔴 Critical System Access</h3>
                                        <p style="margin: 0; color: #dc2626; font-size: 16px;">Monitor system health, manage infrastructure, and respond to platform-wide incidents.</p>
                                    </td>
                                </tr>
                            </table>

                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #991b1b; font-size: 18px;">👥 Global User Management</h3>
                                        <p style="margin: 0; color: #dc2626; font-size: 16px;">Oversee all users, organizations, and roles across the entire platform ecosystem.</p>
                                    </td>
                                </tr>
                            </table>

                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #991b1b; font-size: 18px;">📊 Platform Analytics</h3>
                                        <p style="margin: 0; color: #dc2626; font-size: 16px;">Access comprehensive metrics, audit logs, and compliance reporting for the entire platform.</p>
                                    </td>
                                </tr>
                            </table>

                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #991b1b; font-size: 18px;">🛡️ Security & Compliance</h3>
                                        <p style="margin: 0; color: #dc2626; font-size: 16px;">Manage security policies, conduct audits, and ensure regulatory compliance across all operations.</p>
                                    </td>
                                </tr>
                            </table>

                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="text-align: center; margin: 40px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="{{.DashboardURL}}" style="display: inline-block; background-color: #ef4444; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600;">Access Platform Admin →</a>
                                    </td>
                                </tr>
                            </table>

                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #991b1b; font-size: 18px;">⚠️ Critical Security Responsibilities</h3>
                                        <p style="margin: 0 0 10px 0; color: #dc2626; font-size: 16px;">
                                            As a Platform Administrator, you have access to sensitive system information and user data. Your actions directly impact the security and stability of the entire Ethos Platform.
                                        </p>
                                        <ul style="margin: 10px 0 0 20px; color: #dc2626; font-size: 16px;">
                                            <li>Never share your admin credentials</li>
                                            <li>Report suspicious activity immediately</li>
                                            <li>Follow the principle of least privilege</li>
                                            <li>Maintain detailed audit logs of all actions</li>
                                        </ul>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td style="padding: 20px; text-align: center; font-size: 12px; color: #666; background-color: #f9fafb;">
                ` + canSpamFooter + `
            </td>
        </tr>
    </table>
</body>
</html>`

// Welcome Moderator Template - 95%+ Email Client Compatibility
const welcomeModeratorTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to Ethos - Community Moderator Access</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }
    </style>
</head>
<body>
    <!-- Outlook-safe table-based layout for 95%+ compatibility -->
    <table width="100%" border="0" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: 0 auto; background-color: #ffffff;">
        <tr>
            <td style="background-color: #f97316; padding: 40px 30px; text-align: center;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td style="text-align: center;">
                            <div style="width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px;">
                                <div style="width: 40px; height: 25px; border-bottom: 3px solid #f97316; border-radius: 2px; margin: 0 auto;"></div>
                            </div>
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Community Moderator</h1>
                            <p style="color: #fed7aa; margin: 10px 0 0 0; opacity: 0.9;">Help maintain a healthy community</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>

        <tr>
            <td style="padding: 40px 30px;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td>
                            <h2 style="color: #1f2937; margin-bottom: 20px; font-size: 24px;">Welcome, {{.Name}}! 🛡️</h2>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                You've been selected as a Community Moderator for the Ethos Platform. Your role is crucial in maintaining a respectful, inclusive, and productive environment for all users.
                            </p>

                            <!-- Moderator Feature 1 -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fff7ed; border: 1px solid #fed7aa; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #9a3412; font-size: 18px;">👁️ Content Review</h3>
                                        <p style="margin: 0; color: #c2410c; font-size: 16px;">Review flagged content, assess violations, and take appropriate moderation actions.</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- Moderator Feature 2 -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fff7ed; border: 1px solid #fed7aa; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #9a3412; font-size: 18px;">⚖️ Fair Enforcement</h3>
                                        <p style="margin: 0; color: #c2410c; font-size: 16px;">Apply community guidelines consistently and fairly to all users and content.</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- Moderator Feature 3 -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fff7ed; border: 1px solid #fed7aa; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #9a3412; font-size: 18px;">💬 Community Support</h3>
                                        <p style="margin: 0; color: #c2410c; font-size: 16px;">Help resolve conflicts, answer user questions, and foster positive interactions.</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- Moderator Feature 4 -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fff7ed; border: 1px solid #fed7aa; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #9a3412; font-size: 18px;">📈 Impact Tracking</h3>
                                        <p style="margin: 0; color: #c2410c; font-size: 16px;">Monitor moderation effectiveness and contribute to community health improvements.</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- CTA Button -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="text-align: center; margin: 40px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="{{.DashboardURL}}" style="display: inline-block; background-color: #f97316; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600;">Access Moderation Tools →</a>
                                    </td>
                                </tr>
                            </table>

            <div style="background-color: #fef3c7; border: 1px solid #f59e0b; border-radius: 8px; padding: 20px; margin: 20px 0;">
                <h3 style="margin: 0 0 10px 0; color: #92400e;">Community Guidelines</h3>
                            <p style="margin: 0 0 10px 0; color: #92400e; font-size: 16px;">
                                As a moderator, you're expected to uphold our community standards. Familiarize yourself with the guidelines and use them as your decision-making framework.
                            </p>
                            <p style="margin: 0; color: #92400e; font-size: 16px;">
                                <a href="#" style="color: #92400e; text-decoration: underline;">View Community Guidelines →</a>
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td style="padding: 20px; text-align: center; font-size: 12px; color: #666; background-color: #f9fafb;">
                ` + canSpamFooter + `
            </td>
        </tr>
    </table>
</body>
</html>`

// Welcome Support Engineer Template - 95%+ Email Client Compatibility
const welcomeSupportEngineerTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to Ethos - Support Engineer Access</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }
    </style>
</head>
<body>
    <!-- Outlook-safe table-based layout for 95%+ compatibility -->
    <table width="100%" border="0" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: 0 auto; background-color: #ffffff;">
        <tr>
            <td style="background-color: #06b6d4; padding: 40px 30px; text-align: center;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td style="text-align: center;">
                            <div style="width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px;">
                                <div style="width: 40px; height: 25px; border-bottom: 3px solid #06b6d4; border-radius: 2px; margin: 0 auto;"></div>
                            </div>
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Support Engineer</h1>
                            <p style="color: #a5f3fc; margin: 10px 0 0 0; opacity: 0.9;">Technical support and incident resolution</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>

        <tr>
            <td style="padding: 40px 30px;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td>
                            <h2 style="color: #1f2937; margin-bottom: 20px; font-size: 24px;">Welcome, {{.Name}}! 🔧</h2>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                You've been granted Support Engineer access to the Ethos Platform. Your expertise will help maintain system reliability and resolve technical issues for our users.
                            </p>

                            <!-- Support Feature 1 -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #ecfeff; border: 1px solid #a5f3fc; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #0e7490; font-size: 18px;">🚨 Incident Response</h3>
                                        <p style="margin: 0; color: #0891b2; font-size: 16px;">Respond to system alerts, troubleshoot issues, and coordinate incident resolution.</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- Support Feature 2 -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #ecfeff; border: 1px solid #a5f3fc; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #0e7490; font-size: 18px;">🔍 Diagnostic Tools</h3>
                                        <p style="margin: 0; color: #0891b2; font-size: 16px;">Access advanced monitoring tools, logs, and diagnostic utilities to identify root causes.</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- Support Feature 3 -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #ecfeff; border: 1px solid #a5f3fc; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #0e7490; font-size: 18px;">👤 User Support</h3>
                                        <p style="margin: 0; color: #0891b2; font-size: 16px;">Assist users with technical issues and provide guidance on platform functionality.</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- Support Feature 4 -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #ecfeff; border: 1px solid #a5f3fc; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #0e7490; font-size: 18px;">📊 Performance Monitoring</h3>
                                        <p style="margin: 0; color: #0891b2; font-size: 16px;">Monitor system performance, identify bottlenecks, and optimize platform reliability.</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- CTA Button -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="text-align: center; margin: 40px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="{{.DashboardURL}}" style="display: inline-block; background-color: #06b6d4; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600;">Access Support Dashboard →</a>
                                    </td>
                                </tr>
                            </table>

                            <!-- Support Best Practices -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #ecfdf5; border: 1px solid #6ee7b7; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #065f46; font-size: 18px;">Support Best Practices</h3>
                                        <ul style="margin: 0; color: #047857; padding-left: 20px; font-size: 16px;">
                                            <li>Document all troubleshooting steps and resolutions</li>
                                            <li>Escalate complex issues promptly</li>
                                            <li>Maintain clear communication with stakeholders</li>
                                            <li>Follow incident response protocols</li>
                                        </ul>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td style="padding: 20px; text-align: center; font-size: 12px; color: #666; background-color: #f9fafb;">
                ` + canSpamFooter + `
            </td>
        </tr>
    </table>
</body>
</html>`

// Welcome Compliance Auditor Template - 95%+ Email Client Compatibility
const welcomeComplianceAuditorTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to Ethos - Compliance Auditor Access</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }
    </style>
        .audit-feature p { margin: 0; color: #4338ca; }
    </style>
</head>
<body>
    <!-- Outlook-safe table-based layout for 95%+ compatibility -->
    <table width="100%" border="0" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: 0 auto; background-color: #ffffff;">
        <tr>
            <td style="background-color: #6366f1; padding: 40px 30px; text-align: center;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td style="text-align: center;">
                            <div style="width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px;">
                                <div style="width: 40px; height: 25px; border-bottom: 3px solid #6366f1; border-radius: 2px; margin: 0 auto;"></div>
                            </div>
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Compliance Auditor</h1>
                            <p style="color: #c7d2fe; margin: 10px 0 0 0; opacity: 0.9;">Regulatory compliance and audit oversight</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>

        <tr>
            <td style="padding: 40px 30px;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td>
                            <h2 style="color: #1f2937; margin-bottom: 20px; font-size: 24px;">Welcome, {{.Name}}! 📋</h2>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                You've been granted Compliance Auditor access to the Ethos Platform. Your role ensures our platform maintains the highest standards of regulatory compliance and data protection.
                            </p>

                            <!-- Audit Feature 1 -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #eef2ff; border: 1px solid #c7d2fe; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #312e81; font-size: 18px;">📊 Audit Trail Review</h3>
                                        <p style="margin: 0; color: #3730a3; font-size: 16px;">Access comprehensive audit logs and review system activities for compliance verification.</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- Audit Feature 2 -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #eef2ff; border: 1px solid #c7d2fe; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #312e81; font-size: 18px;">⚖️ Regulatory Compliance</h3>
                                        <p style="margin: 0; color: #3730a3; font-size: 16px;">Monitor adherence to GDPR, SOC2, ISO27001, and other regulatory requirements.</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- Audit Feature 3 -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #eef2ff; border: 1px solid #c7d2fe; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #312e81; font-size: 18px;">🔒 Data Protection</h3>
                                        <p style="margin: 0; color: #3730a3; font-size: 16px;">Verify data handling practices, privacy controls, and security measures.</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- Audit Feature 4 -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #eef2ff; border: 1px solid #c7d2fe; border-radius: 8px; margin: 15px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #312e81; font-size: 18px;">📈 Compliance Reporting</h3>
                                        <p style="margin: 0; color: #3730a3; font-size: 16px;">Generate detailed compliance reports and track remediation of identified issues.</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- CTA Button -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="text-align: center; margin: 40px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="{{.DashboardURL}}" style="display: inline-block; background-color: #6366f1; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600;">Access Audit Dashboard →</a>
                                    </td>
                                </tr>
                            </table>

                            <!-- Audit Independence Notice -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fef3c7; border: 1px solid #f59e0b; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #92400e; font-size: 18px;">Audit Independence</h3>
                                        <p style="margin: 0; color: #92400e; font-size: 16px;">
                                            As a Compliance Auditor, maintain independence and objectivity in all audit activities. Report findings directly to appropriate oversight bodies.
                                        </p>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td style="padding: 20px; text-align: center; font-size: 12px; color: #666; background-color: #f9fafb;">
                ` + canSpamFooter + `
            </td>
        </tr>
    </table>
</body>
</html>`

// Feedback Submitted Template
const feedbackSubmittedTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Feedback Submitted Successfully - Ethos</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }
    </style>
</head>
<body>
    <!-- Outlook-safe table-based layout -->
    <table class="container" width="100%" border="0" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: 0 auto; background-color: #ffffff;">
        <tr>
            <td class="header" style="background-color: #10b981; padding: 40px 30px; text-align: center;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td style="text-align: center;">
                            <div class="logo" style="width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px;">
                                <div style="width: 40px; height: 25px; border-bottom: 3px solid #10b981; border-radius: 2px; margin: 0 auto;"></div>
                            </div>
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Feedback Received</h1>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>

        <tr>
            <td class="content" style="padding: 40px 30px;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td>
                            <h2 style="color: #1f2937; margin-bottom: 20px; font-size: 24px;">Thank you, {{.Name}}! ✅</h2>

                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px; text-align: center;">
                                        <h3 style="margin: 0 0 10px 0; color: #065f46; font-size: 18px;">Your feedback has been successfully submitted</h3>
                                        <p style="margin: 0; color: #047857;">We appreciate you taking the time to share your thoughts with us.</p>
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 20px; font-size: 16px;">
                                Your feedback titled "<strong>{{.FeedbackTitle}}</strong>" has been received and will be reviewed by our team.
                            </p>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                You can track the status of your feedback and engage with the community discussion using the link below:
                            </p>

                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="text-align: center; margin: 30px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="{{.FeedbackURL}}" style="display: inline-block; background-color: #10b981; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600;">View Your Feedback →</a>
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #6b7280; font-size: 14px; margin-bottom: 20px;">
                                Feedback ID: <code style="background-color: #f3f4f6; padding: 2px 6px; border-radius: 4px;">{{.FeedbackID}}</code>
                            </p>

                            <p style="color: #4b5563; line-height: 1.6; font-size: 16px;">
                                We'll notify you of any updates or responses to your feedback. Thank you for helping us improve the Ethos Platform!
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td style="padding: 20px; text-align: center; font-size: 12px; color: #666; background-color: #f9fafb;">
                ` + canSpamFooter + `
            </td>
        </tr>
    </table>
</body>
</html>`

// Appeal Submitted Template
const appealSubmittedTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Appeal Submitted - Ethos Platform</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }
    </style>
</head>
<body>
    <!-- Outlook-safe table-based layout for 95%+ compatibility -->
    <table width="100%" border="0" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: 0 auto; background-color: #ffffff;">
        <tr>
            <td style="background-color: #f59e0b; padding: 40px 30px; text-align: center;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td style="text-align: center;">
                            <div style="width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px;">
                                <div style="width: 40px; height: 25px; border-bottom: 3px solid #f59e0b; border-radius: 2px; margin: 0 auto;"></div>
                            </div>
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Appeal Submitted</h1>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>

        <tr>
            <td style="padding: 40px 30px;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td>
                            <h2 style="color: #1f2937; margin-bottom: 20px; font-size: 24px;">Appeal Received, {{.Name}}</h2>

                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fffbeb; border: 1px solid #fed7aa; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #92400e; font-size: 18px;">Status: Under Review</h3>
                                        <p style="margin: 0; color: #92400e; font-size: 16px;">Your appeal has been successfully submitted and is now being reviewed by our team.</p>
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 20px; font-size: 16px;">
                                We've received your appeal regarding the recent platform decision. Our moderation team will carefully review your submission and respond within 3-5 business days.
                            </p>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                You can check the status of your appeal and receive updates through your dashboard:
                            </p>

                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="text-align: center; margin: 30px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="{{.AppealURL}}" style="display: inline-block; background-color: #f59e0b; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600;">Track Appeal Status →</a>
                                    </td>
                                </tr>
                            </table>

                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #f3f4f6; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h4 style="margin: 0 0 10px 0; color: #374151; font-size: 18px;">What happens next?</h4>
                                        <ol style="margin: 0; color: #4b5563; padding-left: 20px; font-size: 16px;">
                                            <li>Our team reviews your appeal within 24-48 hours</li>
                                            <li>You'll receive an email notification when a decision is made</li>
                                            <li>If approved, changes take effect immediately</li>
                                            <li>If additional information is needed, we'll contact you directly</li>
                                        </ol>
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #6b7280; font-size: 14px; margin-bottom: 20px;">
                                Appeal ID: <code style="background-color: #f3f4f6; padding: 2px 6px; border-radius: 4px;">{{.AppealID}}</code>
                            </p>

                            <p style="color: #6b7280; font-size: 14px;">
                                If you have any questions about your appeal, please don't hesitate to contact our support team.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td style="padding: 20px; text-align: center; font-size: 12px; color: #666; background-color: #f9fafb;">
                ` + canSpamFooter + `
            </td>
        </tr>
    </table>
</body>
</html>`

// Appeal Status Update Template
const appealStatusUpdateTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Appeal Status Update - Ethos Platform</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }
    </style>
</head>
<body>
    <!-- Outlook-safe table-based layout for 95%+ compatibility -->
    <table width="100%" border="0" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: 0 auto; background-color: #ffffff;">
        <tr>
            <td style="background-color: #3b82f6; padding: 40px 30px; text-align: center;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td style="text-align: center;">
                            <div style="width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px;">
                                <div style="width: 40px; height: 25px; border-bottom: 3px solid #3b82f6; border-radius: 2px; margin: 0 auto;"></div>
                            </div>
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Appeal Decision</h1>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>

        <tr>
            <td style="padding: 40px 30px;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td>
                            <h2 style="color: #1f2937; margin-bottom: 20px; font-size: 24px;">Appeal Status Update</h2>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                Dear {{.Name}}, our moderation team has reviewed your appeal and updated its status.
                            </p>

                            <!-- Appeal Information -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #f3f4f6; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h4 style="margin: 0 0 10px 0; color: #374151; font-size: 18px;">Appeal Information:</h4>
                                        <p style="margin: 0 0 5px 0; color: #4b5563; font-size: 16px;"><strong>Appeal ID:</strong> {{.AppealID}}</p>
                                        <p style="margin: 0 0 5px 0; color: #4b5563; font-size: 16px;"><strong>Status:</strong> {{.NewStatus}}</p>
                                        <p style="margin: 0; color: #4b5563; font-size: 16px;"><strong>Original Decision:</strong> {{.OriginalDecision}}</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- Status Update -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #eff6ff; border: 1px solid #bfdbfe; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #1e40af; font-size: 18px;">Status Update: {{.NewStatus}}</h3>
                                        <p style="margin: 0; color: #1e40af; font-size: 16px;">{{.StatusDescription}}</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- Appeal Specifics -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #f9fafb; border-radius: 8px; padding: 20px; margin: 20px 0; border-left: 4px solid #3b82f6;">
                                <tr>
                                    <td>
                                        <h4 style="margin: 0 0 10px 0; color: #374151; font-size: 18px;">Appeal Specifics:</h4>
                                        <p style="margin: 0 0 5px 0; color: #4b5563; font-size: 16px;"><strong>Submitted:</strong> {{.AppealDate}}</p>
                                        <p style="margin: 0 0 5px 0; color: #4b5563; font-size: 16px;"><strong>Last Updated:</strong> {{.LastUpdate}}</p>
                                        <p style="margin: 0 0 5px 0; color: #4b5563; font-size: 16px;"><strong>Reviewed by:</strong> {{.ModeratorName}}</p>
                                        <p style="margin: 0; color: #4b5563; font-size: 16px;"><strong>Estimated Resolution:</strong> {{.EstimatedResolution}}</p>
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                For more details about this decision and next steps, please visit your appeal dashboard:
                            </p>

                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="text-align: center; margin: 30px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="{{.AppealURL}}" style="display: inline-block; background-color: #3b82f6; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600;">View Appeal Details →</a>
                                    </td>
                                </tr>
                            </table>

                            {{if eq .Status "approved"}}
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h4 style="margin: 0 0 10px 0; color: #065f46; font-size: 18px;">What happens next?</h4>
                                        <p style="margin: 0; color: #047857; font-size: 16px;">
                                            The appealed action has been reversed. You should see the changes reflected in your account within the next few minutes.
                                        </p>
                                    </td>
                                </tr>
                            </table>
                            {{end}}

                            {{if eq .Status "rejected"}}
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fef3c7; border: 1px solid #f59e0b; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h4 style="margin: 0 0 10px 0; color: #92400e; font-size: 18px;">Appeal Process</h4>
                                        <p style="margin: 0; color: #92400e; font-size: 16px;">
                                            If you disagree with this decision, you can contact our support team for further assistance.
                                        </p>
                                    </td>
                                </tr>
                            </table>
                            {{end}}

                            <p style="color: #6b7280; font-size: 14px;">
                                Thank you for your patience during the appeal review process.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td style="padding: 20px; text-align: center; font-size: 12px; color: #666; background-color: #f9fafb;">
                ` + canSpamFooter + `
            </td>
        </tr>
    </table>
</body>
</html>`

// Organization Invitation Template
const orgInvitationTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Organization Invitation - Ethos Platform</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }
    </style>
</head>
<body>
    <!-- Outlook-safe table-based layout for 95%+ compatibility -->
    <table width="100%" border="0" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: 0 auto; background-color: #ffffff;">
        <tr>
            <td style="background-color: #8b5cf6; padding: 40px 30px; text-align: center;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td style="text-align: center;">
                            <div style="width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px;">
                                <div style="width: 40px; height: 25px; border-bottom: 3px solid #8b5cf6; border-radius: 2px; margin: 0 auto;"></div>
                            </div>
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Organization Invitation</h1>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>

        <tr>
            <td style="padding: 40px 30px;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td>
                            <h2 style="color: #1f2937; margin-bottom: 20px; font-size: 24px;">You're Invited to Join an Organization</h2>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                Hi {{.Name}}, you've been invited to join the <strong>{{.Organization}}</strong> organization on the Ethos Platform.
                            </p>

                            <!-- Invitation Details -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #faf5ff; border: 1px solid #d8b4fe; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #6b21a8; font-size: 18px;">Invitation Details</h3>
                                        <p style="margin: 0 0 10px 0; color: #7c3aed; font-size: 16px;"><strong>Organization:</strong> {{.Organization}}</p>
                                        <p style="margin: 0 0 10px 0; color: #7c3aed; font-size: 16px;"><strong>Invited by:</strong> {{.InviterName}} ({{.InviterEmail}})</p>
                                        <p style="margin: 0 0 10px 0; color: #7c3aed; font-size: 16px;"><strong>Role:</strong> {{.Role}}</p>
                                        <p style="margin: 0; color: #7c3aed; font-size: 16px;"><strong>Expires:</strong> {{.ExpiryDate}}</p>
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                As a <strong>{{.Role}}</strong>, you'll have access to organization-specific features, team collaboration tools, and shared resources.
                            </p>

                            <!-- CTA Buttons -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="text-align: center; margin: 30px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="{{.AcceptURL}}" style="display: inline-block; background-color: #8b5cf6; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600; margin: 10px;">Accept Invitation</a>
                                        <a href="{{.DeclineURL}}" style="display: inline-block; background-color: #6b7280; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600; margin: 10px;">Decline Invitation</a>
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #6b7280; font-size: 12px; text-align: center; margin: 10px 0;">
                                This invitation will expire on {{.ExpiryDate}}. No immediate action required.
                            </p>

                            <!-- Organization Benefits -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #f3f4f6; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h4 style="margin: 0 0 10px 0; color: #374151; font-size: 18px;">What you'll get:</h4>
                                        <ul style="margin: 0; color: #4b5563; padding-left: 20px; font-size: 16px;">
                                            <li>Access to organization-specific content and discussions</li>
                                            <li>Team collaboration and project management tools</li>
                                            <li>Organization-wide analytics and insights</li>
                                            <li>Customized notifications and preferences</li>
                                        </ul>
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #6b7280; font-size: 14px;">
                                This invitation will expire in 7 days. If you have any questions, please contact the person who invited you.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td style="padding: 20px; text-align: center; font-size: 12px; color: #666; background-color: #f9fafb;">
                ` + canSpamFooter + `
            </td>
        </tr>
    </table>
</body>
</html>`

// Escalation Received Template - 95%+ Email Client Compatibility
const escalationReceivedTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>New Escalation Requires Attention - Ethos</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }
    </style>
</head>
<body>
    <!-- Outlook-safe table-based layout for 95%+ compatibility -->
    <table width="100%" border="0" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: 0 auto; background-color: #ffffff;">
        <tr>
            <td style="background-color: #dc2626; padding: 40px 30px; text-align: center;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td style="text-align: center;">
                            <div style="width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px;">
                                <div style="width: 40px; height: 25px; border-bottom: 3px solid #dc2626; border-radius: 2px; margin: 0 auto;"></div>
                            </div>
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Urgent Escalation</h1>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>

        <tr>
            <td style="padding: 40px 30px;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td>
                            <h2 style="color: #1f2937; margin-bottom: 20px; font-size: 24px;">Escalation Requires Immediate Attention</h2>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                Dear {{.Name}}, a critical issue has been escalated and requires your immediate attention.
                            </p>

                            <!-- Urgent Alert -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #991b1b; font-size: 18px;">Priority: {{.Priority}}</h3>
                                        <p style="margin: 0; color: #dc2626; font-size: 16px;">{{.Description}}</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- Escalation Details -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #f3f4f6; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h4 style="margin: 0 0 10px 0; color: #374151; font-size: 18px;">Recipient Information:</h4>
                                        <p style="margin: 0 0 15px 0; color: #4b5563; font-size: 16px;"><strong>Name:</strong> {{.Name}}</p>

                                        <h4 style="margin: 20px 0 10px 0; color: #374151; font-size: 18px;">Escalation Details:</h4>
                                        <p style="margin: 0 0 5px 0; color: #4b5563; font-size: 16px;"><strong>Escalation ID:</strong> {{.AlertID}}</p>
                                        <p style="margin: 0 0 5px 0; color: #4b5563; font-size: 16px;"><strong>Affected User:</strong> {{.AffectedUser}}</p>
                                        <p style="margin: 0 0 5px 0; color: #4b5563; font-size: 16px;"><strong>Priority:</strong> {{.Priority}}</p>
                                        <p style="margin: 0 0 10px 0; color: #4b5563; font-size: 16px;"><strong>Severity:</strong> {{.Severity}}</p>

                                        <!-- Issue Description -->
                                        <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #ffffff; border-radius: 6px; margin-top: 10px;">
                                            <tr>
                                                <td style="padding: 15px;">
                                                    <p style="margin: 0 0 8px 0; color: #1e40af; font-weight: 600; font-size: 16px;">Issue Description:</p>
                                                    <p style="margin: 0; color: #374151; font-size: 16px;">{{.Description}}</p>
                                                </td>
                                            </tr>
                                        </table>
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                Please review and respond to this escalation as soon as possible to minimize impact on the platform and users.
                            </p>

                            <!-- CTA Button -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="text-align: center; margin: 30px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="{{.EscalationURL}}" style="display: inline-block; background-color: #dc2626; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600;">Review Escalation →</a>
                                    </td>
                                </tr>
                            </table>

                            <!-- Response Guidelines -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fef3c7; border: 1px solid #f59e0b; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h4 style="margin: 0 0 10px 0; color: #92400e; font-size: 18px;">Response Guidelines:</h4>
                                        <ul style="margin: 0; color: #92400e; padding-left: 20px; font-size: 16px;">
                                            <li>Acknowledge receipt within 30 minutes</li>
                                            <li>Provide initial assessment within 2 hours</li>
                                            <li>Escalate further if needed following protocol</li>
                                            <li>Document all actions and communications</li>
                                        </ul>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td style="padding: 20px; text-align: center; font-size: 12px; color: #666; background-color: #f9fafb;">
                ` + canSpamFooter + `
            </td>
        </tr>
    </table>
</body>
</html>`

// Moderation Alert Template - 95%+ Email Client Compatibility
const moderationAlertTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Content Requires Moderation Review - Ethos</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }
    </style>
</head>
<body>
    <!-- Outlook-safe table-based layout for 95%+ compatibility -->
    <table width="100%" border="0" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: 0 auto; background-color: #ffffff;">
        <tr>
            <td style="background-color: #f97316; padding: 40px 30px; text-align: center;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td style="text-align: center;">
                            <div style="width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px;">
                                <div style="width: 40px; height: 25px; border-bottom: 3px solid #f97316; border-radius: 2px; margin: 0 auto;"></div>
                            </div>
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Moderation Review</h1>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>

        <tr>
            <td style="padding: 40px 30px;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td>
                            <h2 style="color: #1f2937; margin-bottom: 20px; font-size: 24px;">Content Requires Your Review</h2>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                Dear {{.Name}}, content on the platform has been flagged for potential policy violations and requires your immediate review.
                            </p>

                            <!-- Alert Box -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fff7ed; border: 1px solid #fed7aa; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #9a3412; font-size: 18px;">Content Flagged for Review</h3>
                                        <p style="margin: 0; color: #c2410c; font-size: 16px;">{{.ContentType}} content has been flagged and requires moderation attention.</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- Content Details -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #f3f4f6; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h4 style="margin: 0 0 10px 0; color: #374151; font-size: 18px;">Content Details:</h4>
                                        <p style="margin: 0 0 5px 0; color: #4b5563; font-size: 16px;"><strong>Content ID:</strong> {{.ContentID}}</p>
                                        <p style="margin: 0 0 5px 0; color: #4b5563; font-size: 16px;"><strong>Type:</strong> {{.ContentType}}</p>
                                        <p style="margin: 0 0 5px 0; color: #4b5563; font-size: 16px;"><strong>Title:</strong> {{.ContentTitle}}</p>
                                        <p style="margin: 0 0 10px 0; color: #4b5563; font-size: 16px;"><strong>Reason:</strong> {{.ReportReason}}</p>

                                        <!-- Content Snippet -->
                                        <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #ffffff; border-radius: 6px; margin-top: 10px;">
                                            <tr>
                                                <td style="padding: 15px;">
                                                    <p style="margin: 0 0 8px 0; color: #1e40af; font-weight: 600; font-size: 16px;">Content Snippet:</p>
                                                    <p style="margin: 0; color: #374151; font-style: italic; font-size: 16px;">"{{.ContentSnippet}}"</p>
                                                </td>
                                            </tr>
                                        </table>
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                {{.ActionRequired}}
                            </p>

                            <!-- CTA Button -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="text-align: center; margin: 30px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="{{.ContentURL}}" style="display: inline-block; background-color: #f97316; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600;">Review Content →</a>
                                    </td>
                                </tr>
                            </table>

                            <!-- Moderation Guidelines -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fef3c7; border: 1px solid #f59e0b; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h4 style="margin: 0 0 10px 0; color: #92400e; font-size: 18px;">Moderation Guidelines:</h4>
                                        <ul style="margin: 0; color: #92400e; padding-left: 20px; font-size: 16px;">
                                            <li>Review content against community guidelines</li>
                                            <li>Document your reasoning for any actions taken</li>
                                            <li>Consider context and intent when making decisions</li>
                                            <li>Escalate complex cases to senior moderators</li>
                                        </ul>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td style="padding: 20px; text-align: center; font-size: 12px; color: #666; background-color: #f9fafb;">
                ` + canSpamFooter + `
            </td>
        </tr>
    </table>
</body>
</html>`

// Audit Report Available Template
const auditReportAvailableTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Audit Report Available - Ethos Platform</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }
    </style>
</head>
<body>
    <!-- Outlook-safe table-based layout for 95%+ compatibility -->
    <table width="100%" border="0" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: 0 auto; background-color: #ffffff;">
        <tr>
            <td style="background-color: #6366f1; padding: 40px 30px; text-align: center;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td style="text-align: center;">
                            <div style="width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px;">
                                <div style="width: 40px; height: 25px; border-bottom: 3px solid #6366f1; border-radius: 2px; margin: 0 auto;"></div>
                            </div>
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Audit Report Ready</h1>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>

        <tr>
            <td style="padding: 40px 30px;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td>
                            <h2 style="color: #1f2937; margin-bottom: 20px; font-size: 24px;">Audit Report Available for Review</h2>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                Dear {{.Name}}, the {{.ReportType}} has been completed and is now available for your review.
                            </p>

                            <!-- Report Card -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #eef2ff; border: 1px solid #c7d2fe; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #312e81; font-size: 18px;">{{.ReportType}} Report</h3>
                                        <p style="margin: 0 0 5px 0; color: #4338ca; font-size: 16px;"><strong>Report ID:</strong> {{.ReportID}}</p>
                                        <p style="margin: 0 0 5px 0; color: #4338ca; font-size: 16px;"><strong>Date Range:</strong> {{.ReportPeriod}}</p>
                                        <p style="margin: 0 0 10px 0; color: #4338ca; font-size: 16px;"><strong>Generated Date:</strong> {{.GeneratedDate}}</p>
                                        <p style="margin: 0 0 10px 0; color: #4338ca; font-size: 16px;"><strong>Findings:</strong></p>

                                        <!-- Key Findings -->
                                        <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #ffffff; border-radius: 6px; margin-top: 10px;">
                                            <tr>
                                                <td style="padding: 15px;">
                                                    <p style="margin: 0 0 8px 0; color: #1e40af; font-weight: 600; font-size: 16px;">Key Findings Summary:</p>
                                                    <p style="margin: 0 0 8px 0; color: #374151; font-size: 16px;">• {{.KeyFindings}}</p>
                                                    <p style="margin: 0 0 8px 0; color: #374151; font-size: 16px;"><strong>Risk Assessment:</strong> {{.RiskAssessment}}</p>
                                                    <p style="margin: 0; color: #374151; font-size: 16px;"><strong>Compliance Status:</strong> {{.ComplianceStatus}}</p>
                                                </td>
                                            </tr>
                                        </table>
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                This report contains {{.FindingsCount}} findings from the audit period. Please review the details and take any necessary follow-up actions.
                            </p>

                            <!-- CTA Button -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="text-align: center; margin: 30px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="{{.ReportURL}}" style="display: inline-block; background-color: #6366f1; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600;">View Audit Report →</a>
                                    </td>
                                </tr>
                            </table>

                            <!-- Report Contents -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h4 style="margin: 0 0 10px 0; color: #065f46; font-size: 18px;">Report Contents:</h4>
                                        <ul style="margin: 0; color: #047857; padding-left: 20px; font-size: 16px;">
                                            <li>Executive summary of audit findings</li>
                                            <li>Detailed analysis of compliance issues</li>
                                            <li>Recommended remediation actions</li>
                                            <li>Supporting evidence and documentation</li>
                                        </ul>
                                    </td>
                                </tr>
                            </table>

                            <!-- Important Notes -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fef3c7; border: 1px solid #f59e0b; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h4 style="margin: 0 0 10px 0; color: #92400e; font-size: 18px;">Important Notes:</h4>
                                        <p style="margin: 0; color: #92400e; font-size: 16px;">
                                            This report contains sensitive compliance information. Please ensure it is stored securely and shared only with authorized personnel.
                                        </p>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td style="padding: 20px; text-align: center; font-size: 12px; color: #666; background-color: #f9fafb;">
                ` + canSpamFooter + `
            </td>
        </tr>
    </table>
</body>
</html>`

// System Health Alert Template
const systemHealthAlertTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>System Health Alert - Ethos Platform</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }
        .container { max-width: 600px; margin: 0 auto; background-color: #ffffff; }
        .header { background: linear-gradient(135deg, #dc2626 0%, #b91c1c 100%); padding: 40px 30px; text-align: center; }
        .logo { width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; display: flex; align-items: center; justify-content: center; }
        .logo div { width: 40px; height: 25px; border-bottom: 3px solid #dc2626; border-radius: 2px; }
        .content { padding: 40px 30px; }
        .button { display: inline-block; background: linear-gradient(135deg, #dc2626 0%, #b91c1c 100%); color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600; margin: 20px 0; }
        .critical { background-color: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; padding: 20px; margin: 20px 0; }
        .warning { background-color: #fffbeb; border: 1px solid #fed7aa; }
        .info { background-color: #eff6ff; border: 1px solid #bfdbfe; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">
                <div></div>
            </div>
            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">System Alert</h1>
        </div>

        <div class="content">
            <h2 style="color: #1f2937; margin-bottom: 20px;">System Health Alert</h2>

            <div class="critical {{if eq .Severity "warning"}}warning{{else if eq .Severity "info"}}info{{end}}">
                <h3 style="margin: 0 0 10px 0; color: #991b1b;">{{.Severity}} Alert</h3>
                <p style="margin: 0; color: #dc2626;">{{.AlertType}} - {{.Service}}</p>
            </div>

            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 20px;">
                Hi {{.Name}}, our monitoring systems have detected a {{.Severity}} condition affecting the Ethos Platform.
            </p>

            <div style="background-color: #f3f4f6; border-radius: 8px; padding: 20px; margin: 20px 0;">
                <h4 style="margin: 0 0 10px 0; color: #374151;">Alert Details:</h4>
                <p style="margin: 0 0 5px 0; color: #4b5563;"><strong>Service:</strong> {{.Service}}</p>
                <p style="margin: 0 0 5px 0; color: #4b5563;"><strong>Incident ID:</strong> {{.IncidentID}}</p>
                <p style="margin: 0; color: #4b5563;"><strong>Description:</strong> {{.Description}}</p>
            </div>

            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px;">
                Our engineering team has been notified and is actively investigating this issue. We'll provide updates as more information becomes available.
            </p>

            <div style="text-align: center; margin: 30px 0;">
                <a href="{{.StatusURL}}" class="button">View System Status →</a>
            </div>

            {{if eq .Severity "critical"}}
            <div style="background-color: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; padding: 20px; margin: 20px 0;">
                <h4 style="margin: 0 0 10px 0; color: #991b1b;">Impact & Response:</h4>
                <p style="margin: 0; color: #dc2626;">
                    This is a critical system issue that may affect platform availability. Our team is working urgently to resolve this matter.
                </p>
            </div>
            {{end}}

            <p style="color: #6b7280; font-size: 14px;">
                For real-time updates, please monitor our status page or subscribe to platform notifications.
            </p>
        </div>

        ` + canSpamFooter + `
    </div>
</body>
</html>`

// Account Deletion Template
const accountDeletionTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Account Deletion Confirmation - Ethos Platform</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }
        .container { max-width: 600px; margin: 0 auto; background-color: #ffffff; }
        .header { background: linear-gradient(135deg, #6b7280 0%, #4b5563 100%); padding: 40px 30px; text-align: center; }
        .logo { width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; display: flex; align-items: center; justify-content: center; }
        .logo div { width: 40px; height: 25px; border-bottom: 3px solid #6b7280; border-radius: 2px; }
        .content { padding: 40px 30px; }
        .warning { background-color: #fef3c7; border: 1px solid #f59e0b; border-radius: 8px; padding: 20px; margin: 20px 0; }
        .info { background-color: #eff6ff; border: 1px solid #bfdbfe; border-radius: 8px; padding: 20px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">
                <div></div>
            </div>
            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Account Deletion</h1>
        </div>

        <div class="content">
            <h2 style="color: #1f2937; margin-bottom: 20px;">Account Deletion Confirmation</h2>

            <div class="warning">
                <h3 style="margin: 0 0 10px 0; color: #92400e;">⚠️ Important Notice</h3>
                <p style="margin: 0; color: #92400e;">
                    Your Ethos Platform account deletion request has been processed successfully.
                </p>
            </div>

            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 20px;">
                Hi {{.Name}}, we're writing to confirm that your Ethos Platform account has been permanently deleted.
            </p>

            <div class="info">
                <h4 style="margin: 0 0 10px 0; color: #1e40af;">What has been deleted:</h4>
                <ul style="margin: 0; color: #1e40af; padding-left: 20px;">
                    <li>Your account profile and personal information</li>
                    <li>All feedback and content you created</li>
                    <li>Your notification preferences and settings</li>
                    <li>Access to any organizations you were a member of</li>
                </ul>
            </div>

            <div style="background-color: #f3f4f6; border-radius: 8px; padding: 20px; margin: 20px 0;">
                <h4 style="margin: 0 0 10px 0; color: #374151;">Data Retention Notice:</h4>
                <p style="margin: 0 0 10px 0; color: #4b5563;">
                    As required by law, some anonymized data may be retained for regulatory compliance purposes. This data cannot be used to identify you personally.
                </p>
                <p style="margin: 0; color: #4b5563;">
                    For more information about our data retention policies, please refer to our Privacy Policy.
                </p>
            </div>

            <div style="background-color: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; padding: 20px; margin: 20px 0;">
                <h4 style="margin: 0 0 10px 0; color: #991b1b;">Account Recovery:</h4>
                <p style="margin: 0; color: #dc2626;">
                    Account deletion is permanent and cannot be undone. If you wish to use Ethos Platform again in the future, you will need to create a new account.
                </p>
            </div>

            <p style="color: #4b5563; line-height: 1.6;">
                Thank you for being part of the Ethos community. We hope our paths cross again in the future.
            </p>
        </div>

        ` + canSpamFooter + `
    </div>
</body>
</html>`

// Security Alert Template
const securityAlertTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Security Alert - Ethos Platform</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }
    </style>
</head>
<body>
    <!-- Outlook-safe table-based layout for 95%+ compatibility -->
    <table width="100%" border="0" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: 0 auto; background-color: #ffffff;">
        <tr>
            <td style="background-color: #dc2626; padding: 40px 30px; text-align: center;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td style="text-align: center;">
                            <div style="width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px;">
                                <div style="width: 40px; height: 25px; border-bottom: 3px solid #dc2626; border-radius: 2px; margin: 0 auto;"></div>
                            </div>
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Security Alert</h1>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>

        <tr>
            <td style="padding: 40px 30px;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td>
                            <h2 style="color: #1f2937; margin-bottom: 20px; font-size: 24px;">Security Event Detected</h2>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                Hi {{.Name}}, our security systems have detected a potential security event related to your Ethos Platform account.
                            </p>

                            <!-- Security Alert Box -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #991b1b; font-size: 18px;">⚠️ Security Event: {{.EventType}}</h3>
                                        <p style="margin: 0; color: #dc2626; font-size: 16px;">
                                            We detected unusual activity on your account that requires your attention.
                                        </p>
                                    </td>
                                </tr>
                            </table>

                            <!-- Event Details -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #f3f4f6; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h4 style="margin: 0 0 10px 0; color: #374151; font-size: 18px;">Event Details:</h4>
                                        <p style="margin: 0 0 5px 0; color: #4b5563; font-size: 16px;"><strong>Event:</strong> {{.EventType}}</p>
                                        <p style="margin: 0 0 5px 0; color: #4b5563; font-size: 16px;"><strong>Time:</strong> {{.EventTime}}</p>
                                        <p style="margin: 0 0 5px 0; color: #4b5563; font-size: 16px;"><strong>IP Address:</strong> {{.IPAddress}}</p>
                                        <p style="margin: 0 0 5px 0; color: #4b5563; font-size: 16px;"><strong>Location:</strong> {{.Location}}</p>
                                        <p style="margin: 0; color: #4b5563; font-size: 16px;"><strong>User Agent:</strong> {{.UserAgent}}</p>
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;">
                                If this activity was not performed by you, please take immediate action to secure your account.
                            </p>

                            <!-- CTA Button -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="text-align: center; margin: 30px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="{{.ActionURL}}" style="display: inline-block; background-color: #dc2626; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600;">Review Account Security →</a>
                                    </td>
                                </tr>
                            </table>

                            <!-- Security Recommendations -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fef3c7; border: 1px solid #f59e0b; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h4 style="margin: 0 0 10px 0; color: #92400e; font-size: 18px;">Security Recommendations:</h4>
                                        <ul style="margin: 0; color: #92400e; padding-left: 20px; font-size: 16px;">
                                            <li>Change your password immediately</li>
                                            <li>Enable two-factor authentication if not already enabled</li>
                                            <li>Review your recent account activity</li>
                                            <li>Contact support if you suspect unauthorized access</li>
                                        </ul>
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #6b7280; font-size: 14px;">
                                This alert was generated automatically by our security monitoring systems. If you have any questions, please contact our support team.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td style="padding: 20px; text-align: center; font-size: 12px; color: #666; background-color: #f9fafb;">
                ` + canSpamFooter + `
            </td>
        </tr>
    </table>
</body>
</html>`



// Feedback Moderated Template
const feedbackModeratedTemplate = "<!DOCTYPE html>\n" +
	"<html lang=\"en\">\n" +
	"<head>\n" +
	"    <meta charset=\"UTF-8\">\n" +
	"    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n" +
	"    <title>Feedback Moderated - Ethos Platform</title>\n" +
	"    <style>\n" +
	"        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }\n" +
	"        .container { max-width: 600px; margin: 0 auto; background-color: #ffffff; }\n" +
	"        .header { background-color: #667eea; padding: 40px 30px; text-align: center; }\n" +
	"        .logo { width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px; }\n" +
	"        .logo div { width: 40px; height: 25px; border-bottom: 3px solid #667eea; border-radius: 2px; margin: 0 auto; }\n" +
	"        .content { padding: 40px 30px; }\n" +
	"        .button { display: inline-block; background-color: #667eea; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600; margin: 20px 0; }\n" +
	"        .status { background-color: #f0f9ff; border: 1px solid #0ea5e9; border-radius: 8px; padding: 20px; margin: 20px 0; }\n" +
	"        .footer { background-color: #f9fafb; padding: 20px 30px; text-align: center; font-size: 14px; color: #6b7280; }\n" +
	"    </style>\n" +
	"</head>\n" +
	"<body>\n" +
	"    <!-- Outlook-safe table-based layout for 95%+ compatibility -->\n" +
	"    <table class=\"container\" width=\"100%\" border=\"0\" cellspacing=\"0\" cellpadding=\"0\" style=\"max-width: 600px; margin: 0 auto; background-color: #ffffff;\">\n" +
	"        <tr>\n" +
	"            <td class=\"header\" style=\"background-color: #667eea; padding: 40px 30px; text-align: center;\">\n" +
	"                <table width=\"100%\" border=\"0\" cellspacing=\"0\" cellpadding=\"0\">\n" +
	"                    <tr>\n" +
	"                        <td style=\"text-align: center;\">\n" +
	"                            <div class=\"logo\" style=\"width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px;\">\n" +
	"                                <div style=\"width: 40px; height: 25px; border-bottom: 3px solid #667eea; border-radius: 2px; margin: 0 auto;\"></div>\n" +
	"                            </div>\n" +
	"                            <h1 style=\"color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;\">Feedback Moderated</h1>\n" +
	"                        </td>\n" +
	"                    </tr>\n" +
	"                </table>\n" +
	"            </td>\n" +
	"        </tr>\n" +
	"\n" +
	"        <tr>\n" +
	"            <td class=\"content\" style=\"padding: 40px 30px;\">\n" +
	"                <table width=\"100%\" border=\"0\" cellspacing=\"0\" cellpadding=\"0\">\n" +
	"                    <tr>\n" +
	"                        <td>\n" +
	"                            <h2 style=\"color: #1f2937; margin-bottom: 20px; font-size: 24px;\">Feedback Review Complete</h2>\n" +
	"\n" +
	"                            <p style=\"color: #4b5563; line-height: 1.6; margin-bottom: 20px; font-size: 16px;\">\n" +
	"                                Dear {{.Name}}, our moderation team has reviewed your feedback submission and taken appropriate action.\n" +
	"                            </p>\n" +
	"\n" +
	"                            <table width=\"100%\" border=\"0\" cellspacing=\"0\" cellpadding=\"0\" style=\"background-color: #f0f9ff; border: 1px solid #0ea5e9; border-radius: 8px; margin: 20px 0;\">\n" +
	"                                <tr>\n" +
	"                                    <td style=\"padding: 20px;\">\n" +
	"                                        <h3 style=\"margin: 0 0 10px 0; color: #0ea5e9; font-size: 18px;\">Feedback Details:</h3>\n" +
	"                                        <p style=\"margin: 0 0 5px 0; color: #374151; font-size: 16px;\"><strong>Title:</strong> {{.FeedbackTitle}}</p>\n" +
	"                                        <p style=\"margin: 0 0 5px 0; color: #374151; font-size: 16px;\"><strong>ID:</strong> {{.FeedbackID}}</p>\n" +
	"                                        <p style=\"margin: 0 0 5px 0; color: #374151; font-size: 16px;\"><strong>Decision:</strong> {{.ModerationReason}}</p>\n" +
	"                                        <p style=\"margin: 0; color: #374151; font-size: 16px;\"><strong>Reviewed by:</strong> {{.ModeratorName}}</p>\n" +
	"                                    </td>\n" +
	"                                </tr>\n" +
	"                            </table>\n" +
	"\n" +
	"                            <p style=\"color: #4b5563; line-height: 1.6; margin-bottom: 30px; font-size: 16px;\">\n" +
	"                                Thank you for contributing to the Ethos Platform community. Your feedback helps us improve our services and maintain a positive environment for all users.\n" +
	"                            </p>\n" +
	"\n" +
	"                            <table width=\"100%\" border=\"0\" cellspacing=\"0\" cellpadding=\"0\" style=\"text-align: center; margin: 30px 0;\">\n" +
	"                                <tr>\n" +
	"                                    <td style=\"text-align: center;\">\n" +
	"                                        <a href=\"{{.DashboardURL}}\" style=\"display: inline-block; background-color: #667eea; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600;\">View Feedback Status →</a>\n" +
	"                                    </td>\n" +
	"                                </tr>\n" +
	"                            </table>\n" +
	"\n" +
	"                            <p style=\"color: #6b7280; font-size: 14px;\">\n" +
	"                                If you have any questions about this moderation decision or would like to submit additional feedback, please don't hesitate to contact our support team.\n" +
	"                            </p>\n" +
	"                        </td>\n" +
	"                    </tr>\n" +
	"                </table>\n" +
	"            </td>\n" +
	"        </tr>\n" +
	"        <tr>\n" +
	"            <td style=\"padding: 20px; text-align: center; font-size: 12px; color: #666; background-color: #f9fafb;\">\n" +
	"                " + canSpamFooter + "\n" +
	"            </td>\n" +
	"        </tr>\n" +
	"    </table>\n" +
	"</body>\n" +
	"</html>"

// Appeal Resolved Template
// Appeal Resolved Template - 95%+ Email Client Compatibility
const appealResolvedTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Appeal Resolution - Ethos Platform</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }
    </style>
</head>
<body>
    <!-- Outlook-safe table-based layout for 95%+ compatibility -->
    <table width="100%" border="0" cellspacing="0" cellpadding="0" style="max-width: 600px; margin: 0 auto; background-color: #ffffff;">
        <tr>
            <td style="background-color: #667eea; padding: 40px 30px; text-align: center;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td style="text-align: center;">
                            <div style="width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; text-align: center; line-height: 60px;">
                                <div style="width: 40px; height: 25px; border-bottom: 3px solid #667eea; border-radius: 2px; margin: 0 auto;"></div>
                            </div>
                            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Appeal Resolution</h1>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>

        <tr>
            <td style="padding: 40px 30px;">
                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                    <tr>
                        <td>
                            <h2 style="color: #1f2937; margin-bottom: 20px; font-size: 24px;">Appeal Decision Finalized</h2>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 20px; font-size: 16px;">
                                Dear {{.Name}}, we have completed our review of your appeal and reached a final decision.
                            </p>

                            <!-- Final Decision -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #f0fdf4; border: 1px solid #22c55e; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h3 style="margin: 0 0 10px 0; color: #15803d; font-size: 18px;">Final Decision: {{.FinalDecision}}</h3>
                                        <p style="margin: 0 0 10px 0; color: #374151; font-size: 16px;"><strong>Appeal ID:</strong> {{.AppealID}}</p>
                                        <p style="margin: 0 0 10px 0; color: #374151; font-size: 16px;"><strong>Resolution Date:</strong> {{.ResolutionDate}}</p>
                                        <p style="margin: 0; color: #374151; font-size: 16px;"><strong>Decision Summary:</strong></p>
                                        <p style="margin: 10px 0 0 0; color: #374151; font-size: 16px;">{{.ResolutionDetails}}</p>
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #4b5563; line-height: 1.6; margin-bottom: 20px; font-size: 16px;">
                                This decision was made after careful consideration of your appeal, relevant platform policies, and community guidelines. Our goal is to maintain a fair and positive environment for all users.
                            </p>

                            <!-- Appeal Timeline -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #f3f4f6; border-radius: 8px; margin: 20px 0;">
                                <tr>
                                    <td style="padding: 20px;">
                                        <h4 style="margin: 0 0 10px 0; color: #374151; font-size: 18px;">Appeal Timeline:</h4>
                                        <p style="margin: 0 0 5px 0; color: #6b7280; font-size: 16px;"><strong>Submitted:</strong> {{.AppealDate}}</p>
                                        <p style="margin: 0 0 5px 0; color: #6b7280; font-size: 16px;"><strong>Reviewed by:</strong> {{.ModeratorName}}</p>
                                        <p style="margin: 0; color: #6b7280; font-size: 16px;"><strong>Resolution:</strong> {{.ResolutionDate}}</p>
                                    </td>
                                </tr>
                            </table>

                            <!-- CTA Button -->
                            <table width="100%" border="0" cellspacing="0" cellpadding="0" style="text-align: center; margin: 30px 0;">
                                <tr>
                                    <td style="text-align: center;">
                                        <a href="{{.DashboardURL}}" style="display: inline-block; background-color: #667eea; color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600;">View Appeal Details →</a>
                                    </td>
                                </tr>
                            </table>

                            <p style="color: #6b7280; font-size: 14px;">
                                If you have any questions about this decision or need further clarification, please contact our support team. For urgent matters, you may also submit a new appeal if circumstances have changed.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td style="padding: 20px; text-align: center; font-size: 12px; color: #666; background-color: #f9fafb;">
                ` + canSpamFooter + `
            </td>
        </tr>
    </table>
</body>
</html>`

// Organization Member Added Template
const orgMemberAddedTemplate = "<!DOCTYPE html>\n" +
	"<html lang=\"en\">\n" +
	"<head>\n" +
	"    <meta charset=\"UTF-8\">\n" +
	"    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n" +
	"    <title>Welcome to {{.Organization}} - Ethos Platform</title>\n" +
	"    <style>\n" +
	"        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 0; background-color: #f9fafb; }\n" +
	"        .container { max-width: 600px; margin: 0 auto; background-color: #ffffff; }\n" +
	"        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); padding: 40px 30px; text-align: center; }\n" +
	"        .logo { width: 60px; height: 60px; background-color: #ffffff; border-radius: 50%; margin: 0 auto 20px; display: flex; align-items: center; justify-content: center; }\n" +
	"        .logo div { width: 40px; height: 25px; border-bottom: 3px solid #667eea; border-radius: 2px; }\n" +
	"        .content { padding: 40px 30px; }\n" +
	"        .button { display: inline-block; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: #ffffff; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600; margin: 20px 0; }\n" +
	"        .org-info { background-color: #f0f9ff; border: 1px solid #0ea5e9; border-radius: 8px; padding: 20px; margin: 20px 0; }\n" +
	"        .footer { background-color: #f9fafb; padding: 20px 30px; text-align: center; font-size: 14px; color: #6b7280; }\n" +
	"    </style>\n" +
	"</head>\n" +
	"<body>\n" +
	"    <div class=\"container\">\n" +
	"        <div class=\"header\">\n" +
	"            <div class=\"logo\">\n" +
	"                <div></div>\n" +
	"            </div>\n" +
	"            <h1 style=\"color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;\">Welcome to {{.Organization}}</h1>\n" +
	"        </div>\n" +
	"\n" +
	"        <div class=\"content\">\n" +
	"            <h2 style=\"color: #1f2937; margin-bottom: 20px;\">Organization Membership Confirmed</h2>\n" +
	"\n" +
	"            <p style=\"color: #4b5563; line-height: 1.6; margin-bottom: 20px;\">\n" +
	"                Congratulations, {{.Name}}! You have been successfully added to {{.Organization}} on the Ethos Platform.\n" +
	"            </p>\n" +
	"\n" +
	"            <div class=\"org-info\">\n" +
	"                <h3 style=\"margin: 0 0 10px 0; color: #0ea5e9;\">Membership Details:</h3>\n" +
	"                <p style=\"margin: 0 0 5px 0; color: #374151;\"><strong>Organization:</strong> {{.Organization}}</p>\n" +
	"                <p style=\"margin: 0 0 5px 0; color: #374151;\"><strong>Role:</strong> {{.MemberRole}}</p>\n" +
	"                <p style=\"margin: 0 0 5px 0; color: #374151;\"><strong>Added by:</strong> {{.AdminName}} ({{.AdminEmail}})</p>\n" +
	"                <p style=\"margin: 0; color: #374151;\"><strong>Join Date:</strong> {{.JoinDate}}</p>\n" +
	"            </div>\n" +
	"\n" +
	"            <p style=\"color: #4b5563; line-height: 1.6; margin-bottom: 20px;\">\n" +
	"                As a member of {{.Organization}}, you now have access to organization-specific features, resources, and collaboration tools on the Ethos Platform. Your role as {{.MemberRole}} provides you with appropriate permissions and access levels.\n" +
	"            </p>\n" +
	"\n" +
	"            <div style=\"background-color: #f3f4f6; border-radius: 8px; padding: 20px; margin: 20px 0;\">\n" +
	"                <h4 style=\"margin: 0 0 10px 0; color: #374151;\">Getting Started:</h4>\n" +
	"                <ul style=\"margin: 0; color: #6b7280; padding-left: 20px;\">\n" +
	"                    <li>Explore your organization's dashboard and resources</li>\n" +
	"                    <li>Connect with other team members</li>\n" +
	"                    <li>Access organization-specific tools and documentation</li>\n" +
	"                    <li>Review your role permissions and responsibilities</li>\n" +
	"                </ul>\n" +
	"            </div>\n" +
	"\n" +
	"            <div style=\"text-align: center; margin: 30px 0;\">\n" +
	"                <a href=\"{{.DashboardURL}}\" class=\"button\">Access {{.Organization}} Dashboard →</a>\n" +
	"            </div>\n" +
	"\n" +
	"            <p style=\"color: #6b7280; font-size: 14px;\">\n" +
	"                If you have any questions about your membership or need assistance getting started, please contact {{.AdminName}} or reach out to our support team.\n" +
	"            </p>\n" +
	"        </div>\n" +
	"\n" +
	"        " + canSpamFooter + "\n" +
	"    </div>\n" +
	"</body>\n" +
	"</html>"
