USE [master]
GO
/****** Object:  Database [DigitalLibrary]    Script Date: 22.02.2021 22:51:24 ******/
CREATE DATABASE [DigitalLibrary]
 CONTAINMENT = NONE
 ON  PRIMARY 
( NAME = N'DigitalLibrary', FILENAME = N'C:\Program Files\Microsoft SQL Server\MSSQL14.MSSQLSERVER\MSSQL\DATA\DigitalLibrary.mdf' , SIZE = 5120KB , MAXSIZE = UNLIMITED, FILEGROWTH = 1024KB )
 LOG ON 
( NAME = N'DigitalLibrary_log', FILENAME = N'C:\Program Files\Microsoft SQL Server\MSSQL14.MSSQLSERVER\MSSQL\DATA\DigitalLibrary_log.ldf' , SIZE = 5184KB , MAXSIZE = 2048GB , FILEGROWTH = 10%)
GO
ALTER DATABASE [DigitalLibrary] SET COMPATIBILITY_LEVEL = 120
GO
IF (1 = FULLTEXTSERVICEPROPERTY('IsFullTextInstalled'))
begin
EXEC [DigitalLibrary].[dbo].[sp_fulltext_database] @action = 'enable'
end
GO
ALTER DATABASE [DigitalLibrary] SET ANSI_NULL_DEFAULT OFF 
GO
ALTER DATABASE [DigitalLibrary] SET ANSI_NULLS OFF 
GO
ALTER DATABASE [DigitalLibrary] SET ANSI_PADDING OFF 
GO
ALTER DATABASE [DigitalLibrary] SET ANSI_WARNINGS OFF 
GO
ALTER DATABASE [DigitalLibrary] SET ARITHABORT OFF 
GO
ALTER DATABASE [DigitalLibrary] SET AUTO_CLOSE OFF 
GO
ALTER DATABASE [DigitalLibrary] SET AUTO_SHRINK OFF 
GO
ALTER DATABASE [DigitalLibrary] SET AUTO_UPDATE_STATISTICS ON 
GO
ALTER DATABASE [DigitalLibrary] SET CURSOR_CLOSE_ON_COMMIT OFF 
GO
ALTER DATABASE [DigitalLibrary] SET CURSOR_DEFAULT  GLOBAL 
GO
ALTER DATABASE [DigitalLibrary] SET CONCAT_NULL_YIELDS_NULL OFF 
GO
ALTER DATABASE [DigitalLibrary] SET NUMERIC_ROUNDABORT OFF 
GO
ALTER DATABASE [DigitalLibrary] SET QUOTED_IDENTIFIER OFF 
GO
ALTER DATABASE [DigitalLibrary] SET RECURSIVE_TRIGGERS OFF 
GO
ALTER DATABASE [DigitalLibrary] SET  DISABLE_BROKER 
GO
ALTER DATABASE [DigitalLibrary] SET AUTO_UPDATE_STATISTICS_ASYNC OFF 
GO
ALTER DATABASE [DigitalLibrary] SET DATE_CORRELATION_OPTIMIZATION OFF 
GO
ALTER DATABASE [DigitalLibrary] SET TRUSTWORTHY OFF 
GO
ALTER DATABASE [DigitalLibrary] SET ALLOW_SNAPSHOT_ISOLATION OFF 
GO
ALTER DATABASE [DigitalLibrary] SET PARAMETERIZATION SIMPLE 
GO
ALTER DATABASE [DigitalLibrary] SET READ_COMMITTED_SNAPSHOT OFF 
GO
ALTER DATABASE [DigitalLibrary] SET HONOR_BROKER_PRIORITY OFF 
GO
ALTER DATABASE [DigitalLibrary] SET RECOVERY FULL 
GO
ALTER DATABASE [DigitalLibrary] SET  MULTI_USER 
GO
ALTER DATABASE [DigitalLibrary] SET PAGE_VERIFY CHECKSUM  
GO
ALTER DATABASE [DigitalLibrary] SET DB_CHAINING OFF 
GO
ALTER DATABASE [DigitalLibrary] SET FILESTREAM( NON_TRANSACTED_ACCESS = OFF ) 
GO
ALTER DATABASE [DigitalLibrary] SET TARGET_RECOVERY_TIME = 0 SECONDS 
GO
ALTER DATABASE [DigitalLibrary] SET DELAYED_DURABILITY = DISABLED 
GO
EXEC sys.sp_db_vardecimal_storage_format N'DigitalLibrary', N'ON'
GO
ALTER DATABASE [DigitalLibrary] SET QUERY_STORE = OFF
GO
USE [DigitalLibrary]
GO
/****** Object:  User [user]    Script Date: 22.02.2021 22:51:24 ******/
CREATE USER [user] FOR LOGIN [user] WITH DEFAULT_SCHEMA=[dbo]
GO
/****** Object:  User [stas]    Script Date: 22.02.2021 22:51:24 ******/
CREATE USER [stas] WITHOUT LOGIN WITH DEFAULT_SCHEMA=[dbo]
GO
ALTER ROLE [db_owner] ADD MEMBER [user]
GO
/****** Object:  Table [dbo].[Administrators]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Administrators](
	[IDAdmin] [int] NOT NULL,
	[IDUser] [int] NOT NULL,
 CONSTRAINT [PK_Administrators] PRIMARY KEY CLUSTERED 
(
	[IDAdmin] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[Authors]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Authors](
	[IDAuthor] [int] IDENTITY(1,1) NOT NULL,
	[NameAuthor] [varchar](100) NOT NULL,
	[DescribeAuthor] [text] NOT NULL,
	[IDUser] [int] NULL,
 CONSTRAINT [PK_Authors] PRIMARY KEY CLUSTERED 
(
	[IDAuthor] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY] TEXTIMAGE_ON [PRIMARY]
GO
/****** Object:  Table [dbo].[BookAuthor]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[BookAuthor](
	[IDBook] [int] NOT NULL,
	[IDAuthor] [int] NOT NULL,
 CONSTRAINT [PK_BookAuthor] PRIMARY KEY CLUSTERED 
(
	[IDBook] ASC,
	[IDAuthor] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[Books]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Books](
	[IDBook] [int] IDENTITY(1,1) NOT NULL,
	[CoverBook] [image] NULL,
	[NameBook] [varchar](350) NOT NULL,
	[IDGenre] [int] NOT NULL,
	[IDTypeBook] [int] NULL,
	[DescribeBook] [text] NOT NULL,
	[PriceBook] [money] NULL,
 CONSTRAINT [PK_Books] PRIMARY KEY CLUSTERED 
(
	[IDBook] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY] TEXTIMAGE_ON [PRIMARY]
GO
/****** Object:  Table [dbo].[Cart]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Cart](
	[IDCart] [int] NOT NULL,
	[IDUser] [int] NOT NULL,
	[IDBook] [int] NOT NULL,
	[IDCartStatus] [int] NOT NULL,
 CONSTRAINT [PK_Cart] PRIMARY KEY CLUSTERED 
(
	[IDCart] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[CartStatuses]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[CartStatuses](
	[IDCartStatus] [int] NOT NULL,
	[CartStatus] [varchar](20) NOT NULL,
 CONSTRAINT [PK_CartStatuses] PRIMARY KEY CLUSTERED 
(
	[IDCartStatus] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[Formats]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Formats](
	[IDFormat] [int] IDENTITY(1,1) NOT NULL,
	[NameFormat] [varchar](10) NOT NULL,
 CONSTRAINT [PK_Formats] PRIMARY KEY CLUSTERED 
(
	[IDFormat] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[Formulars]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Formulars](
	[IDForm] [int] IDENTITY(1,1) NOT NULL,
	[IDCard] [int] NOT NULL,
	[IDBook] [int] NOT NULL,
	[DateIssue] [date] NOT NULL,
	[DateCompletion] [date] NOT NULL,
	[Returned] [bit] NOT NULL,
 CONSTRAINT [PK_Formulars] PRIMARY KEY CLUSTERED 
(
	[IDForm] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[Genres]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Genres](
	[IDGenre] [int] IDENTITY(1,1) NOT NULL,
	[NameGenre] [varchar](50) NOT NULL,
	[IDParentGenre] [int] NULL,
	[DescribeGenre] [text] NOT NULL,
 CONSTRAINT [PK_Genres] PRIMARY KEY CLUSTERED 
(
	[IDGenre] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY] TEXTIMAGE_ON [PRIMARY]
GO
/****** Object:  Table [dbo].[Partners]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Partners](
	[IDPartnership] [int] NOT NULL,
	[IDUser] [int] NOT NULL,
	[AnnotationText] [text] NOT NULL,
	[CoverText] [image] NOT NULL,
	[Text] [varbinary](max) NOT NULL,
	[IDPartnershipStatus] [int] NOT NULL,
 CONSTRAINT [PK_Partners] PRIMARY KEY CLUSTERED 
(
	[IDPartnership] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY] TEXTIMAGE_ON [PRIMARY]
GO
/****** Object:  Table [dbo].[PartnershipStatuses]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[PartnershipStatuses](
	[IDPartnershipStatus] [int] NOT NULL,
	[PartnershipStatus] [varchar](20) NOT NULL,
 CONSTRAINT [PK_PartnershipStatuses] PRIMARY KEY CLUSTERED 
(
	[IDPartnershipStatus] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[Readers]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Readers](
	[IDReader] [int] IDENTITY(1,1) NOT NULL,
	[Login] [varchar](20) NOT NULL,
	[Name] [varchar](50) NOT NULL,
	[Surname] [varchar](50) NOT NULL,
	[Age] [int] NOT NULL,
 CONSTRAINT [PK_Readers] PRIMARY KEY CLUSTERED 
(
	[IDReader] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[ReadersCards]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[ReadersCards](
	[IDCard] [int] IDENTITY(1,1) NOT NULL,
	[IDUser] [int] NOT NULL,
	[DateCard] [date] NOT NULL,
 CONSTRAINT [PK_ReadersCards] PRIMARY KEY CLUSTERED 
(
	[IDCard] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[TypeBooks]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[TypeBooks](
	[IDTypeBook] [int] IDENTITY(1,1) NOT NULL,
	[NameTypeBook] [varchar](10) NOT NULL,
	[IDFormat] [int] NOT NULL,
 CONSTRAINT [PK_TypeBooks] PRIMARY KEY CLUSTERED 
(
	[IDTypeBook] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[UserCollection]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[UserCollection](
	[IDUser] [int] NOT NULL,
	[IDBook] [int] NOT NULL,
 CONSTRAINT [PK_UserCollection] PRIMARY KEY CLUSTERED 
(
	[IDUser] ASC,
	[IDBook] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[Users]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Users](
	[IDUser] [int] IDENTITY(1,1) NOT NULL,
	[Nickname] [varchar](50) NOT NULL,
	[Email] [varchar](50) NOT NULL,
	[Hash] [varchar](256) NOT NULL,
 CONSTRAINT [PK_Users] PRIMARY KEY CLUSTERED 
(
	[IDUser] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
ALTER TABLE [dbo].[Administrators]  WITH CHECK ADD  CONSTRAINT [FK_Administrators_Users] FOREIGN KEY([IDUser])
REFERENCES [dbo].[Users] ([IDUser])
GO
ALTER TABLE [dbo].[Administrators] CHECK CONSTRAINT [FK_Administrators_Users]
GO
ALTER TABLE [dbo].[Authors]  WITH CHECK ADD  CONSTRAINT [FK_Authors_Users] FOREIGN KEY([IDUser])
REFERENCES [dbo].[Users] ([IDUser])
GO
ALTER TABLE [dbo].[Authors] CHECK CONSTRAINT [FK_Authors_Users]
GO
ALTER TABLE [dbo].[BookAuthor]  WITH CHECK ADD  CONSTRAINT [FK_BookAuthor_Authors] FOREIGN KEY([IDAuthor])
REFERENCES [dbo].[Authors] ([IDAuthor])
GO
ALTER TABLE [dbo].[BookAuthor] CHECK CONSTRAINT [FK_BookAuthor_Authors]
GO
ALTER TABLE [dbo].[BookAuthor]  WITH CHECK ADD  CONSTRAINT [FK_BookAuthor_Books] FOREIGN KEY([IDBook])
REFERENCES [dbo].[Books] ([IDBook])
GO
ALTER TABLE [dbo].[BookAuthor] CHECK CONSTRAINT [FK_BookAuthor_Books]
GO
ALTER TABLE [dbo].[Books]  WITH CHECK ADD  CONSTRAINT [FK_Books_Genres] FOREIGN KEY([IDGenre])
REFERENCES [dbo].[Genres] ([IDGenre])
GO
ALTER TABLE [dbo].[Books] CHECK CONSTRAINT [FK_Books_Genres]
GO
ALTER TABLE [dbo].[Books]  WITH CHECK ADD  CONSTRAINT [FK_Books_TypeBooks] FOREIGN KEY([IDTypeBook])
REFERENCES [dbo].[TypeBooks] ([IDTypeBook])
GO
ALTER TABLE [dbo].[Books] CHECK CONSTRAINT [FK_Books_TypeBooks]
GO
ALTER TABLE [dbo].[Cart]  WITH CHECK ADD  CONSTRAINT [FK_Cart_Books] FOREIGN KEY([IDBook])
REFERENCES [dbo].[Books] ([IDBook])
GO
ALTER TABLE [dbo].[Cart] CHECK CONSTRAINT [FK_Cart_Books]
GO
ALTER TABLE [dbo].[Cart]  WITH CHECK ADD  CONSTRAINT [FK_Cart_CartStatuses] FOREIGN KEY([IDCartStatus])
REFERENCES [dbo].[CartStatuses] ([IDCartStatus])
GO
ALTER TABLE [dbo].[Cart] CHECK CONSTRAINT [FK_Cart_CartStatuses]
GO
ALTER TABLE [dbo].[Cart]  WITH CHECK ADD  CONSTRAINT [FK_Cart_Users] FOREIGN KEY([IDUser])
REFERENCES [dbo].[Users] ([IDUser])
GO
ALTER TABLE [dbo].[Cart] CHECK CONSTRAINT [FK_Cart_Users]
GO
ALTER TABLE [dbo].[Formulars]  WITH CHECK ADD  CONSTRAINT [FK_Formulars_Books] FOREIGN KEY([IDBook])
REFERENCES [dbo].[Books] ([IDBook])
GO
ALTER TABLE [dbo].[Formulars] CHECK CONSTRAINT [FK_Formulars_Books]
GO
ALTER TABLE [dbo].[Formulars]  WITH CHECK ADD  CONSTRAINT [FK_Formulars_ReadersCards] FOREIGN KEY([IDCard])
REFERENCES [dbo].[ReadersCards] ([IDCard])
GO
ALTER TABLE [dbo].[Formulars] CHECK CONSTRAINT [FK_Formulars_ReadersCards]
GO
ALTER TABLE [dbo].[Genres]  WITH CHECK ADD  CONSTRAINT [FK_Genres_Genres] FOREIGN KEY([IDParentGenre])
REFERENCES [dbo].[Genres] ([IDGenre])
GO
ALTER TABLE [dbo].[Genres] CHECK CONSTRAINT [FK_Genres_Genres]
GO
ALTER TABLE [dbo].[Partners]  WITH CHECK ADD  CONSTRAINT [FK_Partners_PartnershipStatuses] FOREIGN KEY([IDPartnershipStatus])
REFERENCES [dbo].[PartnershipStatuses] ([IDPartnershipStatus])
GO
ALTER TABLE [dbo].[Partners] CHECK CONSTRAINT [FK_Partners_PartnershipStatuses]
GO
ALTER TABLE [dbo].[Partners]  WITH CHECK ADD  CONSTRAINT [FK_Partners_Users] FOREIGN KEY([IDUser])
REFERENCES [dbo].[Users] ([IDUser])
GO
ALTER TABLE [dbo].[Partners] CHECK CONSTRAINT [FK_Partners_Users]
GO
ALTER TABLE [dbo].[ReadersCards]  WITH CHECK ADD  CONSTRAINT [FK_ReadersCards_Readers] FOREIGN KEY([IDUser])
REFERENCES [dbo].[Readers] ([IDReader])
GO
ALTER TABLE [dbo].[ReadersCards] CHECK CONSTRAINT [FK_ReadersCards_Readers]
GO
ALTER TABLE [dbo].[ReadersCards]  WITH CHECK ADD  CONSTRAINT [FK_ReadersCards_Users] FOREIGN KEY([IDUser])
REFERENCES [dbo].[Users] ([IDUser])
GO
ALTER TABLE [dbo].[ReadersCards] CHECK CONSTRAINT [FK_ReadersCards_Users]
GO
ALTER TABLE [dbo].[TypeBooks]  WITH CHECK ADD  CONSTRAINT [FK_TypeBooks_Formats] FOREIGN KEY([IDFormat])
REFERENCES [dbo].[Formats] ([IDFormat])
GO
ALTER TABLE [dbo].[TypeBooks] CHECK CONSTRAINT [FK_TypeBooks_Formats]
GO
ALTER TABLE [dbo].[UserCollection]  WITH CHECK ADD  CONSTRAINT [FK_UserCollection_Books] FOREIGN KEY([IDBook])
REFERENCES [dbo].[Books] ([IDBook])
GO
ALTER TABLE [dbo].[UserCollection] CHECK CONSTRAINT [FK_UserCollection_Books]
GO
ALTER TABLE [dbo].[UserCollection]  WITH CHECK ADD  CONSTRAINT [FK_UserCollection_Users] FOREIGN KEY([IDUser])
REFERENCES [dbo].[Users] ([IDUser])
GO
ALTER TABLE [dbo].[UserCollection] CHECK CONSTRAINT [FK_UserCollection_Users]
GO
/****** Object:  StoredProcedure [dbo].[GetAllGenres]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE PROCEDURE [dbo].[GetAllGenres] 
AS
BEGIN
Select NameGenre,DescribeGenre 
from Genres
END

GO
/****** Object:  StoredProcedure [dbo].[GetAuthor]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE PROCEDURE [dbo].[GetAuthor]
	@NameAuthor varchar(350) OUTPUT
AS
BEGIN
Select NameAuthor, DescribeAuthor
	FROM Books,BookAuthor,Authors,Genres 
	WHERE Books.IDBook = BookAuthor.IDBook
	AND BookAuthor.IDBook = Authors.IDAuthor
	AND Genres.IDGenre = Books.IDGenre
	AND NameAuthor = @NameAuthor
END
GO
/****** Object:  StoredProcedure [dbo].[GetBook]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE PROCEDURE [dbo].[GetBook]
	@NameBook varchar(350) OUTPUT
AS
BEGIN
Select NameBook, DescribeBook, NameAuthor, NameGenre
	FROM Books,BookAuthor,Authors,Genres 
	WHERE Books.IDBook = BookAuthor.IDBook
	AND BookAuthor.IDBook = Authors.IDAuthor
	AND Genres.IDGenre = Books.IDGenre
	AND NameBook = @NameBook
END
GO
/****** Object:  StoredProcedure [dbo].[GetGenre]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE PROCEDURE [dbo].[GetGenre]
	@NameGenre varchar(350) OUTPUT
AS
BEGIN
Select NameGenre,DescribeGenre,NameBook
	FROM Books,BookAuthor,Authors, Genres 
	WHERE Books.IDBook = BookAuthor.IDBook
	AND BookAuthor.IDBook = Authors.IDAuthor
	AND Genres.IDGenre = Books.IDGenre
	AND NameGenre =  @NameGenre 
END
GO
/****** Object:  StoredProcedure [dbo].[GetTopGenres]    Script Date: 22.02.2021 22:51:24 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE PROCEDURE [dbo].[GetTopGenres]
AS
BEGIN
 SELECT TOP 5 NameGenre FROM Genres
END
GO
USE [master]
GO
ALTER DATABASE [DigitalLibrary] SET  READ_WRITE 
GO
